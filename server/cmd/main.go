package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	v1 "github.com/wcygan/link/generated/go/link/v1"
	"github.com/wcygan/link/generated/go/link/v1/linkv1connect"
)

type UrlShortenerServer struct {
}

func (s *UrlShortenerServer) ShortenUrl(
	ctx context.Context,
	req *connect.Request[v1.ShortenUrlRequest],
) (*connect.Response[v1.ShortenUrlResponse], error) {
	log.Println("Request headers: ", req.Header())
	// TODO: Implement actual URL shortening logic
	shortenedURL := "http://short.url/" + req.Msg.OriginalUrl[:5]
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
	// TODO: Implement actual URL retrieval logic
	originalURL := "http://original.url/" + req.Msg.ShortenedUrl[len("http://short.url/"):]
	res := connect.NewResponse(&v1.GetOriginalUrlResponse{
		OriginalUrl: originalURL,
	})
	res.Header().Set("Shortener-Version", "v1")
	return res, nil
}

func main() {
	urlShortener := &UrlShortenerServer{}
	mux := http.NewServeMux()
	path, handler := linkv1connect.NewUrlShortenerServiceHandler(urlShortener)
	mux.Handle(path, handler)

	fmt.Println("Server starting on :8080")
	err := http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
