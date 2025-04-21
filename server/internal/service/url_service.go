package service

import (
	"context"

	"buf.build/gen/go/wcygan/link/connectrpc/go/link/v1/linkv1connect"
	linkv1 "buf.build/gen/go/wcygan/link/protocolbuffers/go/link/v1"
	"connectrpc.com/connect"
)

// UrlServer implements the linkv1.UrlService interface.
type UrlServer struct{}

// ShortenURL implements linkv1connect.UrlServiceHandler.ShortenURL.
func (s *UrlServer) ShortenURL(
	ctx context.Context,
	req *connect.Request[linkv1.ShortenURLRequest],
) (*connect.Response[linkv1.ShortenURLResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, nil)
}

// RefreshPreview implements linkv1connect.UrlServiceHandler.RefreshPreview.
func (s *UrlServer) RefreshPreview(
	ctx context.Context,
	req *connect.Request[linkv1.RefreshPreviewRequest],
) (*connect.Response[linkv1.RefreshPreviewResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, nil)
}

// GetPreview implements linkv1connect.UrlServiceHandler.GetPreview.
func (s *UrlServer) GetPreview(
	ctx context.Context,
	req *connect.Request[linkv1.GetPreviewRequest],
) (*connect.Response[linkv1.GetPreviewResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, nil)
}

// GenerateQRCode implements linkv1connect.UrlServiceHandler.GenerateQRCode.
func (s *UrlServer) GenerateQRCode(
	ctx context.Context,
	req *connect.Request[linkv1.GenerateQRCodeRequest],
) (*connect.Response[linkv1.GenerateQRCodeResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, nil)
}

// Ensure UrlServer implements the interface
var _ linkv1connect.UrlServiceHandler = (*UrlServer)(nil) 