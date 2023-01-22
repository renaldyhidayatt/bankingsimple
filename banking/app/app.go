package app

import (
	"banking/domain"
	"banking/service"
	"bankinglib/logger"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	_ "github.com/lib/pq"

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
	dbClient := getDbClient()
	customerRepositoryDb := domain.NewCustomerRepositoryDb(dbClient)
	accountRepositoryDb := domain.NewAccountRepositoryDb(dbClient)

	ch := CustomerHandlers{service.NewCustomerService(customerRepositoryDb)}
	ah := AccountHandler{service.NewAccountService(accountRepositoryDb)}

	if runtime.NumCPU() > 2 {
		runtime.GOMAXPROCS(runtime.NumCPU() / 2)
	}
	am := AuthMiddleware{domain.NewAuthRepository()}

	router.Use(am.authorizationHandler())
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome"))
	})

	router.Route("/customer", func(r chi.Router) {
		r.Get("/", ch.getAllCustomers)
		r.Get("/{customer_id}", ch.getCustomer)
		r.Post("/{customer_id}/account", ah.NewAccount)
		r.Post("/{customer_id}/account/{account_id}", ah.MakeTransaction)
	})

	address := os.Getenv("SERVER_ADDRESS")
	port := os.Getenv("SERVER_PORT")

	logger.Info(fmt.Sprintf("Starting server on %s:%s ...", address, port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), router))
}

func getDbClient() *sqlx.DB {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")

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
