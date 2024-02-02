package app

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/labstack/echo/v4"
	"github.com/lewisd1996/baozi-zhongwen/service"
)

type App struct {
	id     string
	Router *echo.Echo
	Auth   *service.AuthService
}

func NewApp() *App {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2"),
	})
	if err != nil {
		// Handle session creation error
	}
	svc := cognitoidentityprovider.New(session)
	authService := service.NewAuthService(svc)

	return &App{
		id:     "baozi-zhongwen",
		Router: echo.New(),
		Auth:   authService,
	}
}
