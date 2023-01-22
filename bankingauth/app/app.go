package app

import (
	"bankingauth/domain"
	"bankingauth/service"
	"bankinglib/logger"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func sanityCheck() {
	envProps := []string{
		"DB_USER",
		"DB_PASSWORD",
		"DB_HOST",
		"DB_PORT",
		"DB_NAME",
		"SERVER_ADDRESS",
		"SERVER_PORT",
	}
	for _, k := range envProps {
		if os.Getenv(k) == "" {
			logger.Fatal(fmt.Sprintf("Environment variable %s not defined. Terminating application...", k))
		}
	}
}

func Start() {
	sanityCheck()
	router := chi.NewRouter()

	authRepository := domain.NewAuthRepository(getDbClient())
	ah := AuthHandler{service.NewLoginService(authRepository, domain.GetRolePermissions())}

	router.Route("/auth", func(r chi.Router) {
		r.Post("/login", ah.Login)
		r.Post("/register", ah.NotImplementedHandler)
		r.Post("/refresh", ah.Refresh)
		r.Post("/verify", ah.Verify)
	})
}

func getDbClient() *sqlx.DB {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	fmt.Println(dbHost + " " + dbPort + " " + dbUser + dbName + dbPassword)

	dataSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sqlx.Connect("postgres", dataSource)

	if err != nil {
		log.Fatal(err.Error())
	}

	db.SetMaxIdleConns(10)

	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	return db
}
