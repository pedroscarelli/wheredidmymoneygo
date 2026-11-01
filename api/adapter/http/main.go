package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/pedroscarelli/wheredidmymoneygo/adapter/http/middleware"
	"github.com/pedroscarelli/wheredidmymoneygo/adapter/postgres"
	"github.com/pedroscarelli/wheredidmymoneygo/core/security"
	"github.com/pedroscarelli/wheredidmymoneygo/di"
)

var requiredEnv = []string{"API_PORT", "DB_USER", "DB_PASSWORD", "DB_HOST", "DB_PORT", "DB_NAME", "JWT_SECRET_KEY"}

func main() {
	for _, name := range requiredEnv {
		if os.Getenv(name) == "" {
			log.Fatalf("missing required environment variable %s", name)
		}
	}

	jwtSecret, err := security.JwtSecret()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	conn := postgres.GetConnection(ctx)
	defer conn.Close()

	if err := postgres.RunMigrations(); err != nil {
		log.Fatalf("error running migrations: %v", err)
	}

	accountService, accountUseCase := di.ConfigAccountDI(conn)
	authenticationService := di.ConfigAuthenticationDI(conn)

	router := mux.NewRouter()
	router.HandleFunc("/login", authenticationService.Login).Methods(http.MethodPost)
	router.HandleFunc("/account", accountService.Create).Methods(http.MethodPost)

	protectedRouter := router.PathPrefix("").Subrouter()
	protectedRouter.Use(middleware.JwtAuthorizer(jwtSecret, accountUseCase))

	protectedRouter.HandleFunc("/account", accountService.Update).Methods(http.MethodPatch)
	protectedRouter.HandleFunc("/account/movement", accountService.Movement).Methods(http.MethodPost)
	protectedRouter.HandleFunc("/account/{id}", accountService.GetByID).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/account/{id}", accountService.Delete).Methods(http.MethodDelete)

	server := &http.Server{
		Addr:              ":" + os.Getenv("API_PORT"),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("listening on port %s", os.Getenv("API_PORT"))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
