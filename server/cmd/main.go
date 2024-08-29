package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"crypto/rand"
	"encoding/base64"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	v1 "github.com/wcygan/link/generated/go/link/v1"
	"github.com/wcygan/link/generated/go/link/v1/linkv1connect"
)

type UrlShortenerServer struct {
	db *pgxpool.Pool
}

func (s *UrlShortenerServer) ShortenUrl(
	ctx context.Context,
	req *connect.Request[v1.ShortenUrlRequest],
) (*connect.Response[v1.ShortenUrlResponse], error) {
	log.Println("Request headers: ", req.Header())
	
	shortCode, err := generateShortCode()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to generate short code: %w", err))
	}

	shortenedURL := fmt.Sprintf("http://short.url/%s", shortCode)

	_, err = s.db.Exec(ctx, "INSERT INTO urls (short_code, original_url) VALUES ($1, $2)", shortCode, req.Msg.OriginalUrl)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to insert URL: %w", err))
	}

	res := connect.NewResponse(&v1.ShortenUrlResponse{
		ShortenedUrl: shortenedURL,
	})
	res.Header().Set("Shortener-Version", "v1")
	return res, nil
}

func (s *UrlShortenerServer) GetOriginalUrl(
	ctx context.Context,
	req *connect.Request[v1.GetOriginalUrlRequest],
) (*connect.Response[v1.GetOriginalUrlResponse], error) {
	log.Println("Request headers: ", req.Header())
	
	shortCode := req.Msg.ShortenedUrl[len("http://short.url/"):]

	var originalURL string
	err := s.db.QueryRow(ctx, "SELECT original_url FROM urls WHERE short_code = $1", shortCode).Scan(&originalURL)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("URL not found: %w", err))
	}

	res := connect.NewResponse(&v1.GetOriginalUrlResponse{
		OriginalUrl: originalURL,
	})
	res.Header().Set("Shortener-Version", "v1")
	return res, nil
}

func generateShortCode() (string, error) {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:6], nil
}

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

	urlShortener := &UrlShortenerServer{db: pool}
	mux := http.NewServeMux()
	path, handler := linkv1connect.NewUrlShortenerServiceHandler(urlShortener)
	mux.Handle(path, handler)

	fmt.Println("Server starting on :8080")
	err = http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func verifyDatabaseSchema(pool *pgxpool.Pool) error {
	ctx := context.Background()
	var tableExists bool
	err := pool.QueryRow(ctx, "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'urls')").Scan(&tableExists)
	if err != nil {
		return fmt.Errorf("failed to check if table exists: %w", err)
	}
	if !tableExists {
		return fmt.Errorf("urls table does not exist")
	}
	return nil
}
