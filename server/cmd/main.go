package main

import (
	"log"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"link/server/internal/service" // Path to our service implementation

	"buf.build/gen/go/wcygan/link/connectrpc/go/link/v1/linkv1connect" // Generated handler constructor
)

func main() {
	urlServer := &service.UrlServer{} // Instantiate our service implementation
	mux := http.NewServeMux()

	// The generated NewUrlServiceHandler constructor routes requests to the service.
	path, handler := linkv1connect.NewUrlServiceHandler(urlServer)
	mux.Handle(path, handler)

	log.Println("Starting server on :8080")
	// Use h2c so we can serve HTTP/2 without TLS.
	err := http.ListenAndServe(
		"localhost:8080",
		h2c.NewHandler(mux, &http2.Server{}),
	)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
