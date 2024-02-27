package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/lewisd1996/baozi-zhongwen/internal/dao"
	"github.com/lewisd1996/baozi-zhongwen/internal/util"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

/* ---------------------------------- Types --------------------------------- */

type AuthService struct {
	cognitoClient     *cognitoidentityprovider.CognitoIdentityProvider
	jwkSet            jwk.Set
	clientId          string
	googleOauthConfig *oauth2.Config
	TokenService      *TokenService
}

type OAuthCodeExchangeResponse struct {
	IdToken      string `json:"id_token"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

/* -------------------------------------------------------------------------- */
/*                                    Init                                    */
/* -------------------------------------------------------------------------- */

func NewAuthService(cognitoClient *cognitoidentityprovider.CognitoIdentityProvider) *AuthService {
	cognitoClientId := os.Getenv("AWS_COGNITO_CLIENT_ID")
	cognitoUserPoolId := os.Getenv("AWS_COGNITO_USER_POOL_ID")
	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	if cognitoClientId == "" {
		panic("missing AWS_COGNITO_CLIENT_ID")
	}
	if cognitoUserPoolId == "" {
		panic("missing AWS_COGNITO_USER_POOL_ID")
	}
	if googleClientId == "" {
		panic("missing GOOGLE_CLIENT_ID")
	}
	if googleClientSecret == "" {
		panic("missing GOOGLE_CLIENT_SECRET")
	}

	var googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("URL") + "/auth/oauth/google/callback",
		ClientID:     googleClientId,
		ClientSecret: googleClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	jwkSet := getJWKSet(fmt.Sprintf("https://cognito-idp.eu-west-2.amazonaws.com/%s/.well-known/jwks.json", cognitoUserPoolId))

	TokenService := NewTokenService()

	return &AuthService{
		cognitoClient:     cognitoClient,
		jwkSet:            jwkSet,
		clientId:          cognitoClientId,
		googleOauthConfig: googleOauthConfig,
		TokenService:      TokenService,
	}
}

/* ---------------------------------- Login --------------------------------- */

func (service *AuthService) LoginWithUsernamePassword(c echo.Context, username, password string) error {
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(username),
			"PASSWORD": aws.String(password),
		},
		ClientId: aws.String(service.clientId),
	}

	result, err := service.cognitoClient.InitiateAuth(input)

	if err != nil {
		if err.Error() == "UserNotConfirmedException: User is not confirmed." {
			encodedUsername := url.QueryEscape(username)
			c.Response().Header().Set("HX-Redirect", "/register/confirm?username="+encodedUsername)
			return c.NoContent(http.StatusOK)
		}

		return err
	}

	accessToken := *result.AuthenticationResult.AccessToken
	refreshToken := *result.AuthenticationResult.RefreshToken

	if accessToken == "" || refreshToken == "" {
		println("error getting tokens")
		return echo.NewHTTPError(http.StatusInternalServerError, "error getting tokens")
	}

	service.TokenService.SetAccessTokenCookie(c, accessToken)
	service.TokenService.SetRefreshTokenCookie(c, refreshToken)

	return nil
}

/* --------------------------------- Logout --------------------------------- */

func (service *AuthService) Logout(c echo.Context) error {
	accessTokenCookie, err := c.Cookie("access_token")

	if err != nil {
		println("missing token cookie")
		return echo.NewHTTPError(http.StatusUnauthorized, "missing token cookie")
	}

	accessToken := accessTokenCookie.Value
	service.cognitoClient.GlobalSignOut(&cognitoidentityprovider.GlobalSignOutInput{
		AccessToken: aws.String(accessToken),
	})

	service.TokenService.DeleteAccessTokenCookie(c)
	service.TokenService.DeleteRefreshTokenCookie(c)

	return nil
}

/* -------------------------------- Register -------------------------------- */

func (service *AuthService) RegisterWithUsernameAndPassword(username, password string) (*cognitoidentityprovider.SignUpOutput, error) {
	err := util.ValidatePassword(password)
	if err != nil {
		return nil, err
	}

	input := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(service.clientId),
		Username: aws.String(username),
		Password: aws.String(password),
	}

	result, err := service.cognitoClient.SignUp(input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

/* --------------------------------- Confirm -------------------------------- */

func (service *AuthService) Confirm(username, code string) error {
	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(service.clientId),
		Username:         aws.String(username),
		ConfirmationCode: aws.String(code),
	}

	_, err := service.cognitoClient.ConfirmSignUp(input)
	if err != nil {
		return err
	}
	return nil
}

func (service *AuthService) ResendConfirmationCode(username string) error {
	input := &cognitoidentityprovider.ResendConfirmationCodeInput{
		ClientId: aws.String(service.clientId),
		Username: aws.String(username),
	}

	_, err := service.cognitoClient.ResendConfirmationCode(input)
	if err != nil {
		return err
	}
	return nil
}

/* ---------------------------------- Token --------------------------------- */

func (service *AuthService) ValidateToken(accessToken string) (jwt.Token, error) {
	token, err := jwt.ParseString(accessToken, jwt.WithKeySet(service.jwkSet), jwt.WithAcceptableSkew(1*time.Minute))
	if err != nil {
		println("[ValidateToken]: Error parsing token:", err.Error())
		return nil, err
	}

	return token, nil
}

func (service *AuthService) RefreshToken(refreshToken string) (*cognitoidentityprovider.AuthenticationResultType, error) {
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("REFRESH_TOKEN_AUTH"),
		AuthParameters: map[string]*string{
			"REFRESH_TOKEN": aws.String(refreshToken),
		},
		ClientId: aws.String(service.clientId),
	}

	result, err := service.cognitoClient.InitiateAuth(input)
	if err != nil {
		return nil, err
	}

	return result.AuthenticationResult, nil
}

func getJWKSet(jwkUrl string) jwk.Set {
	jwkCache := jwk.NewCache(context.Background())
	err := jwkCache.Register(jwkUrl, jwk.WithMinRefreshInterval(10*time.Minute))
	if err != nil {
		panic("failed to register jwk location")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// fetch once on application startup
	_, err = jwkCache.Refresh(ctx, jwkUrl)
	if err != nil {
		panic("failed to fetch on startup")
	}
	// create the cached key set
	return jwk.NewCachedSet(jwkCache, jwkUrl)
}

/* ---------------------------------- OAuth --------------------------------- */

func (service *AuthService) GetGoogleLoginURL() string {
	url := fmt.Sprintf("%s/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&identity_provider=Google", os.Getenv("AWS_COGNITO_URL"), os.Getenv("AWS_COGNITO_CLIENT_ID"), os.Getenv("URL")+"/auth/oauth/google/login/callback")
	return url
}

func (service *AuthService) GetGoogleRegisterURL() string {
	url := fmt.Sprintf("%s/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&identity_provider=Google", os.Getenv("AWS_COGNITO_URL"), os.Getenv("AWS_COGNITO_CLIENT_ID"), os.Getenv("URL")+"/auth/oauth/google/register/callback")
	return url
}

func (service *AuthService) SignInWithGoogle(c echo.Context, code string, dao *dao.Dao) error {
	tokenResponse, err := service.exchangeLoginCodeForToken(code)

	if err != nil {
		println("error exchanging code for token:", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "error exchanging code for token")
	}

	accessToken := tokenResponse.AccessToken
	refreshToken := tokenResponse.RefreshToken

	if accessToken == "" || refreshToken == "" {
		println("error getting tokens")
		return echo.NewHTTPError(http.StatusInternalServerError, "error getting tokens")
	}

	// Get user email from Cognito
	userAuth, err := service.cognitoClient.GetUser(&cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(accessToken),
	})

	if err != nil {
		println("error getting user from cognito:", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "error getting user from cognito")
	}

	// First user attribute is sub
	id := *userAuth.UserAttributes[0].Value

	// Check if user exists in database
	user, err := dao.GetUserById(id)

	if err != nil || user.Email == "" {
		// User does not exist, remove cognito user and return error
		if err.Error() == "qrm: no rows in result set" {
			_, err := service.cognitoClient.DeleteUser(&cognitoidentityprovider.DeleteUserInput{
				AccessToken: aws.String(accessToken),
			})
			if err != nil {
				println("error deleting user from cognito:", err.Error())
			}
			return fmt.Errorf("user is not registered")
		}
		return err
	}

	service.TokenService.SetAccessTokenCookie(c, accessToken)
	service.TokenService.SetRefreshTokenCookie(c, refreshToken)

	return nil
}

func (service *AuthService) RegisterWithGoogle(c echo.Context, code string, dao *dao.Dao) error {
	tokenResponse, err := service.exchangeRegisterCodeForToken(code)

	if err != nil {
		println("error exchanging code for token:", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "error exchanging code for token")
	}

	accessToken := tokenResponse.AccessToken
	refreshToken := tokenResponse.RefreshToken

	if accessToken == "" || refreshToken == "" {
		println("[RegisterWithGoogle] - error getting tokens")
		return echo.NewHTTPError(http.StatusInternalServerError, "error getting tokens")
	}

	// Get user email from Cognito
	userAuth, err := service.cognitoClient.GetUser(&cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(accessToken),
	})

	if err != nil {
		println("[RegisterWithGoogle] - error getting user from cognito:", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "error getting user from cognito")
	}

	for _, attr := range userAuth.UserAttributes {
		println("[RegisterWithGoogle] - user attribute:", *attr.Name, *attr.Value)
	}

	sub := *userAuth.UserAttributes[0].Value
	email := *userAuth.UserAttributes[3].Value

	println("Creating this user in db:", sub, email)

	// Check if user exists in database
	_, err = dao.GetUserByEmail(email)

	if err == nil {
		println("[RegisterWithGoogle] - user already exists")
		_, err := service.cognitoClient.DeleteUser(&cognitoidentityprovider.DeleteUserInput{
			AccessToken: aws.String(accessToken),
		})
		if err != nil {
			println("error deleting user from cognito:", err.Error())
		}
		return fmt.Errorf("user email already exists")
	}

	// Create user in database
	err = dao.CreateUser(email, uuid.MustParse(sub))

	if err != nil {
		println("[RegisterWithGoogle] - error creating user:", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "error creating user")
	}

	println("[RegisterWithGoogle] - Got these tokens:", accessToken, refreshToken)

	service.TokenService.SetAccessTokenCookie(c, accessToken)
	service.TokenService.SetRefreshTokenCookie(c, refreshToken)

	return nil
}

func (service *AuthService) exchangeLoginCodeForToken(code string) (OAuthCodeExchangeResponse, error) {
	client := &http.Client{}
	tokenEndpoint := fmt.Sprintf("%s/oauth2/token", os.Getenv("AWS_COGNITO_URL"))

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", os.Getenv("AWS_COGNITO_CLIENT_ID"))
	data.Set("code", code)
	data.Set("redirect_uri", fmt.Sprintf("%s/auth/oauth/google/login/callback", os.Getenv("URL")))

	req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		println("[AuthService.exchangeLoginCodeForToken]: error creating request:", err.Error())
		return OAuthCodeExchangeResponse{}, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		// Handle error
		println("[AuthService.exchangeLoginCodeForToken]: error making request:", err.Error())
		return OAuthCodeExchangeResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// Handle error
		println("[AuthService.exchangeLoginCodeForToken]: error reading response:", err.Error())
		return OAuthCodeExchangeResponse{}, err
	}

	var tokenResponse OAuthCodeExchangeResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		// Handle error
		println("[AuthService.exchangeLoginCodeForToken]: error parsing response:", err.Error())
		return OAuthCodeExchangeResponse{}, err
	}

	return tokenResponse, nil
}

func (service *AuthService) exchangeRegisterCodeForToken(code string) (OAuthCodeExchangeResponse, error) {
	client := &http.Client{}
	tokenEndpoint := fmt.Sprintf("%s/oauth2/token", os.Getenv("AWS_COGNITO_URL"))

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", os.Getenv("AWS_COGNITO_CLIENT_ID"))
	data.Set("code", code)
	data.Set("redirect_uri", fmt.Sprintf("%s/auth/oauth/google/register/callback", os.Getenv("URL")))

	req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		println("[exchangeRegisterCodeForToken] - error creating request:", err.Error())
		return OAuthCodeExchangeResponse{}, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		// Handle error
		println("[exchangeRegisterCodeForToken] - error making request:", err.Error())
		return OAuthCodeExchangeResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// Handle error
		println("[exchangeRegisterCodeForToken] - error reading response:", err.Error())
		return OAuthCodeExchangeResponse{}, err
	}

	var tokenResponse OAuthCodeExchangeResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		// Handle error
		println("[exchangeRegisterCodeForToken] - error parsing response:", err.Error())
		return OAuthCodeExchangeResponse{}, err
	}

	return tokenResponse, nil
}
