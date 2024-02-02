package service

import (
	"context"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// Initialize a new AuthService

type AuthService struct {
	cognitoClient *cognitoidentityprovider.CognitoIdentityProvider
	jwkSet        jwk.Set
}

func NewAuthService(cognitoClient *cognitoidentityprovider.CognitoIdentityProvider) *AuthService {
	jwkSet := getJWKSet("https://cognito-idp.eu-west-2.amazonaws.com/eu-west-2_1ItHH7zDJ/.well-known/jwks.json")
	return &AuthService{cognitoClient: cognitoClient, jwkSet: jwkSet}
}

// Login

func (service *AuthService) Login(username, password string) (*cognitoidentityprovider.AuthenticationResultType, error) {
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(username),
			"PASSWORD": aws.String(password),
		},
		ClientId: aws.String("5a2vaqjfqsuko38tioiletjh9e"),
	}

	result, err := service.cognitoClient.InitiateAuth(input)
	if err != nil {
		return nil, err
	}
	return result.AuthenticationResult, nil
}

// Logout

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

// Utility functions

func (service *AuthService) ValidateToken(accessToken string) (jwt.Token, error) {
	token, err := jwt.ParseString(accessToken, jwt.WithKeySet(service.jwkSet), jwt.WithAcceptableSkew(1*time.Minute))

	if err != nil {
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
		ClientId: aws.String("5a2vaqjfqsuko38tioiletjh9e"),
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
