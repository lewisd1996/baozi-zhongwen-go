package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/lewisd1996/baozi-zhongwen/util"
)

type AuthService struct {
	cognitoClient *cognitoidentityprovider.CognitoIdentityProvider
	jwkSet        jwk.Set
	clientId      string
}

/* -------------------------------------------------------------------------- */
/*                                    Init                                    */
/* -------------------------------------------------------------------------- */

func NewAuthService(cognitoClient *cognitoidentityprovider.CognitoIdentityProvider) *AuthService {
	cognitoClientId := os.Getenv("AWS_COGNITO_CLIENT_ID")
	cognitoUserPoolId := os.Getenv("AWS_COGNITO_USER_POOL_ID")

	if cognitoClientId == "" {
		panic("missing AWS_COGNITO_CLIENT_ID")
	}
	if cognitoUserPoolId == "" {
		panic("missing AWS_COGNITO_USER_POOL_ID")
	}

	jwkSet := getJWKSet(fmt.Sprintf("https://cognito-idp.eu-west-2.amazonaws.com/%s/.well-known/jwks.json", cognitoUserPoolId))

	return &AuthService{
		cognitoClient: cognitoClient,
		jwkSet:        jwkSet,
		clientId:      cognitoClientId,
	}
}

/* ---------------------------------- Login --------------------------------- */

func (service *AuthService) Login(username, password string) (*cognitoidentityprovider.AuthenticationResultType, error) {
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
		return nil, err
	}
	return result.AuthenticationResult, nil
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

	deletedAccessTokenCookie := http.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}
	c.SetCookie(&deletedAccessTokenCookie)

	deletedRefreshTokenCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}
	c.SetCookie(&deletedRefreshTokenCookie)

	return nil
}

/* -------------------------------- Register -------------------------------- */

func (service *AuthService) Register(username, password string) (*cognitoidentityprovider.SignUpOutput, error) {
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
		println("error parsing token:", err.Error())
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

	// register a minimum refresh interval for this URL.
	// when not specified, defaults to Cache-Control and similar resp headers
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
