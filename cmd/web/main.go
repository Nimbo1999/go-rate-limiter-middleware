package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nimbo1999/go-rate-limit/internal/config"
	"github.com/nimbo1999/go-rate-limit/internal/controllers"
	"github.com/nimbo1999/go-rate-limit/internal/db"
	"github.com/nimbo1999/go-rate-limit/internal/middlewares"
	"github.com/nimbo1999/go-rate-limit/internal/repositories"
	"github.com/nimbo1999/go-rate-limit/internal/services"
	"github.com/spf13/viper"
)

func init() {
	if err := config.LoadEnvConfiguration(); err != nil {
		log.Fatalln(err.Error())
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	log.Println("Server closed")
}

func run() error {
	// Handle SIGINT (CTRL + C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	rdb, err := db.NewRedisClient(context.Background())
	if err != nil {
		return err
	}

	rateLimiterRepository := services.NewRateLimiterRedisChecker(
		repositories.NewRateLimiterRedisRepository(rdb, viper.GetDuration("API_BLOQUER_DURATION_IN_SECONDS")),
		viper.GetInt("IP_MAX_CALLS"),
		viper.GetInt("API_KEY_MAX_CALLS"),
	)

	// Defining the handler
	router := chi.NewRouter()
	// Add all the moddlewares
	router.Use(
		middleware.RealIP,
		middlewares.ApiKeyRetriever,
		middlewares.RateLimiter(rateLimiterRepository),
	)
	router.Handle("/", &controllers.Handler{})
	server := http.Server{
		Addr:    fmt.Sprintf(":%s", viper.GetString("SERVER_PORT")),
		Handler: router,
	}

	srvErr := make(chan error, 1)
	go func() {
		srvErr <- server.ListenAndServe()
	}()

	// Wait for interruption.
	select {
	case err := <-srvErr:
		log.Println("Closing server. There was an error starting the server")
		// Error when starting HTTP server.
		return err
	case <-ctx.Done():
		// Wait for first CTRL+C
		// Stop receiving signal notifications as soon as possible
		log.Println("User has sent a stop signal, closing server...")
		stop()
	}

	// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
	return server.Shutdown(context.Background())
}
