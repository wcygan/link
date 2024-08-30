package server

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	v1 "github.com/wcygan/link/generated/go/auth/v1"
)

type AuthServer struct {
	db *pgxpool.Pool
}

func NewAuthServer(db *pgxpool.Pool) *AuthServer {
	return &AuthServer{db: db}
}

func (s *AuthServer) Register(
	ctx context.Context,
	req *connect.Request[v1.RegisterRequest],
) (*connect.Response[v1.RegisterResponse], error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Msg.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to hash password: %w", err))
	}

	var userID int
	err = s.db.QueryRow(ctx,
		"INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id",
		req.Msg.Username, req.Msg.Email, string(hashedPassword)).Scan(&userID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to insert user: %w", err))
	}

	return connect.NewResponse(&v1.RegisterResponse{
		UserId: fmt.Sprintf("%d", userID),
	}), nil
}

func (s *AuthServer) Login(
	ctx context.Context,
	req *connect.Request[v1.LoginRequest],
) (*connect.Response[v1.LoginResponse], error) {
	var storedHash string
	err := s.db.QueryRow(ctx,
		"SELECT password_hash FROM users WHERE username = $1",
		req.Msg.Username).Scan(&storedHash)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user not found"))
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Msg.Password))
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("invalid password"))
	}

	// TODO: Implement proper token generation
	token := "dummy_token"

	return connect.NewResponse(&v1.LoginResponse{
		Token: token,
	}), nil
}