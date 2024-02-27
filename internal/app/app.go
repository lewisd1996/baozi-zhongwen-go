package app

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/labstack/echo/v4"
	"github.com/lewisd1996/baozi-zhongwen/internal/dao"
	"github.com/lewisd1996/baozi-zhongwen/internal/service"
)

type App struct {
	id     string
	Router *echo.Echo
	Auth   *service.AuthService
	Dao    *dao.Dao
}

func NewApp() *App {
	// ☁️ AWS
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2"),
	})
	if err != nil {
		println("Error creating AWS session")
		panic(err)
	}
	svc := cognitoidentityprovider.New(session)
	authService := service.NewAuthService(svc)

	// 🐘 PostgreSQL
	pgUser := os.Getenv("POSTGRES_USER")
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	pgHost := os.Getenv("POSTGRES_HOST")
	pgPort := os.Getenv("POSTGRES_PORT")
	pgDbName := os.Getenv("POSTGRES_DB")

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", pgUser, pgPassword, pgHost, pgPort, pgDbName)

	db, dbErr := sql.Open("postgres", connectionString)

	if dbErr != nil {
		println("Error opening database connection")
		panic(dbErr)
	}

	if err := db.Ping(); err != nil {
		println("Error pinging database")
		panic(err)
	}

	// 📦 DAO
	dao := dao.NewDao(db)

	return &App{
		id:     "baozi-zhongwen",
		Router: echo.New(),
		Auth:   authService,
		Dao:    dao,
	}
}