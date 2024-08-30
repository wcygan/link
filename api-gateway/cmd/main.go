package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/wcygan/link/api-gateway/internal/gateway"
	authv1connect "github.com/wcygan/link/generated/go/auth/v1/authv1connect"
	linkv1connect "github.com/wcygan/link/generated/go/link/v1/linkv1connect"
)

func main() {
	authClient := authv1connect.NewAuthServiceClient(
		&http.Client{Timeout: 10 * time.Second},
		"http://auth-service:8080",
	)

	linkClient := linkv1connect.NewUrlShortenerServiceClient(
		&http.Client{Timeout: 10 * time.Second},
		"http://link-service:8080",
	)

	gatewayServer := gateway.NewGatewayServer(authClient, linkClient)

	mux := http.NewServeMux()

	mux.Handle(authv1connect.NewAuthServiceHandler(gatewayServer))
	mux.Handle(linkv1connect.NewUrlShortenerServiceHandler(gatewayServer))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: h2c.NewHandler(mux, &http2.Server{}),
		// Set timeouts to ensure connections are properly closed
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Println("API Gateway starting on :8080")
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}