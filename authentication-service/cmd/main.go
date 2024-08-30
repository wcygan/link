package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/wcygan/link/authentication-service/internal/server"
	"github.com/wcygan/link/generated/go/auth/v1/authv1connect"
)

func main() {
	// Get database connection string from environment variable
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// Create a connection pool
	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Failed to parse database URL: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Failed to create connection pool: %v", err)
	}
	defer pool.Close()

	// Verify database connection and schema
	if err := verifyDatabaseSchema(pool); err != nil {
		log.Fatalf("Database schema verification failed: %v", err)
	}

	authServer := server.NewAuthServer(pool)
	mux := http.NewServeMux()
	path, handler := authv1connect.NewAuthServiceHandler(authServer)
	mux.Handle(path, handler)

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Authentication server starting on :8080")
	err = http.ListenAndServe(
		":8080",
		h2c.NewHandler(mux, &http2.Server{}),
	)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func verifyDatabaseSchema(pool *pgxpool.Pool) error {
	ctx := context.Background()
	var tableExists bool
	err := pool.QueryRow(ctx, "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users')").Scan(&tableExists)
	if err != nil {
		return fmt.Errorf("failed to check if table exists: %w", err)
	}
	if !tableExists {
		return fmt.Errorf("users table does not exist")
	}
	return nil
}