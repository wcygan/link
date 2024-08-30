package gateway

import (
	"context"
	"log"

	"connectrpc.com/connect"

	authv1 "github.com/wcygan/link/generated/go/auth/v1"
	authv1connect "github.com/wcygan/link/generated/go/auth/v1/authv1connect"
	linkv1 "github.com/wcygan/link/generated/go/link/v1"
	linkv1connect "github.com/wcygan/link/generated/go/link/v1/linkv1connect"
)

type GatewayServer struct {
	authClient authv1connect.AuthServiceClient
	linkClient linkv1connect.UrlShortenerServiceClient
}

func NewGatewayServer(authClient authv1connect.AuthServiceClient, linkClient linkv1connect.UrlShortenerServiceClient) *GatewayServer {
	return &GatewayServer{
		authClient: authClient,
		linkClient: linkClient,
	}
}

func (s *GatewayServer) Register(ctx context.Context, req *connect.Request[authv1.RegisterRequest]) (*connect.Response[authv1.RegisterResponse], error) {
	resp, err := s.authClient.Register(ctx, req)
	if err != nil {
		log.Printf("Error in Register: %v", err)
	}
	return resp, err
}

func (s *GatewayServer) Login(ctx context.Context, req *connect.Request[authv1.LoginRequest]) (*connect.Response[authv1.LoginResponse], error) {
	resp, err := s.authClient.Login(ctx, req)
	if err != nil {
		log.Printf("Error in Login: %v", err)
	}
	return resp, err
}

func (s *GatewayServer) ShortenUrl(ctx context.Context, req *connect.Request[linkv1.ShortenUrlRequest]) (*connect.Response[linkv1.ShortenUrlResponse], error) {
	resp, err := s.linkClient.ShortenUrl(ctx, req)
	if err != nil {
		log.Printf("Error in ShortenUrl: %v", err)
	}
	return resp, err
}

func (s *GatewayServer) GetOriginalUrl(ctx context.Context, req *connect.Request[linkv1.GetOriginalUrlRequest]) (*connect.Response[linkv1.GetOriginalUrlResponse], error) {
	resp, err := s.linkClient.GetOriginalUrl(ctx, req)
	if err != nil {
		log.Printf("Error in GetOriginalUrl: %v", err)
	}
	return resp, err
}