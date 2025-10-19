package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	httpAdapter "github.com/evgenijurbanovskij/re-partners-assignment/internal/adapters/http"
	"github.com/evgenijurbanovskij/re-partners-assignment/internal/infra/postgres"
	"github.com/evgenijurbanovskij/re-partners-assignment/internal/usecase"
)

const (
	defaultPort    = "8080"
	defaultVersion = "dev"
)

type VersionResponse struct {
	Version string `json:"version"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	version := os.Getenv("VERSION")
	if version == "" {
		version = defaultVersion
	}

	// Initialize components
	// Create slog logger with JSON handler
	slogLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	logger := httpAdapter.NewSlogAdapter(slogLogger)
	solver := usecase.NewDPSolver()

	// Optional PostgreSQL connection
	var db *sqlx.DB
	var dbCleanup func()
	if dbEnabled := os.Getenv("DB_ENABLED"); dbEnabled == "true" {
		log.Println("PostgreSQL integration enabled")

		cfg := postgres.Config{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			Database:        getEnv("DB_NAME", "re_partners"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    25,
			MaxIdleConns:    25,
			ConnMaxLifetime: 5 * time.Minute,
		}

		var err error
		db, err = postgres.Connect(cfg)
		if err != nil {
			log.Printf("Warning: failed to connect to PostgreSQL: %v", err)
			log.Println("Running without database (calculations will not be persisted)")
		} else {
			log.Println("PostgreSQL connected successfully")
			dbCleanup = func() {
				if err := postgres.Close(db); err != nil {
					log.Printf("Error closing database: %v", err)
				}
			}
		}
	} else {
		log.Println("PostgreSQL integration disabled (set DB_ENABLED=true to enable)")
	}

	// Create handler with optional repository
	packHandler := httpAdapter.NewPackHandler(solver, logger)
	if db != nil {
		repo := postgres.NewRepository(db)
		adapter := postgres.NewRepositoryAdapter(repo)
		packHandler = packHandler.WithRepository(adapter)
		log.Println("Database repository integrated with API")
	}

	// Create chi router
	r := chi.NewRouter()

	// Apply middleware
	r.Use(middleware.RequestID)
	r.Use(httpAdapter.RecoveryMiddleware(logger))
	r.Use(httpAdapter.CorrelationIDMiddleware(logger))
	r.Use(httpAdapter.MetricsMiddleware(logger))

	// Health check endpoint
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Version endpoint
	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(VersionResponse{Version: version})
	})

	// Metrics endpoint (Prometheus format)
	r.Handle("/metrics", promhttp.Handler())

	// Pack solver endpoint
	r.Post("/packs/solve", packHandler.SolvePacks)

	// Static files (web UI)
	fs := http.FileServer(http.Dir("./web"))
	r.Handle("/*", fs)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	// Start the server
	go func() {
		log.Printf("Server starting on port %s", port)
		serverErrors <- server.ListenAndServe()
	}()

	// Channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Fatalf("Error starting server: %v", err)

	case sig := <-shutdown:
		log.Printf("Received signal %v, starting graceful shutdown", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Asking listener to shut down and shed load.
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown did not complete in %v: %v", 30*time.Second, err)
			if err := server.Close(); err != nil {
				log.Fatalf("Could not stop server gracefully: %v", err)
			}
		}

		// Close database if connected
		if dbCleanup != nil {
			dbCleanup()
			log.Println("Database connection closed")
		}

		log.Println("Server stopped gracefully")
	}
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
