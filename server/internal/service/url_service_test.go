package service_test

import (
	"context"
	"link/server/internal/service"
	"testing"

	linkv1 "buf.build/gen/go/wcygan/link/protocolbuffers/go/link/v1"
	"connectrpc.com/connect"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUrlService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UrlService Suite")
}

var _ = Describe("UrlServer", func() {
	var (
		urlServer *service.UrlServer
		ctx       context.Context
	)

	BeforeEach(func() {
		urlServer = &service.UrlServer{}
		ctx = context.Background()
	})

	Context("ShortenURL", func() {
		It("should return unimplemented", func() {
			req := connect.NewRequest(&linkv1.ShortenURLRequest{})
			_, err := urlServer.ShortenURL(ctx, req)
			Expect(err).To(HaveOccurred())
			Expect(connect.CodeOf(err)).To(Equal(connect.CodeUnimplemented))
		})
	})

	Context("RefreshPreview", func() {
		It("should return unimplemented", func() {
			req := connect.NewRequest(&linkv1.RefreshPreviewRequest{})
			_, err := urlServer.RefreshPreview(ctx, req)
			Expect(err).To(HaveOccurred())
			Expect(connect.CodeOf(err)).To(Equal(connect.CodeUnimplemented))
		})
	})

	Context("GetPreview", func() {
		It("should return unimplemented", func() {
			req := connect.NewRequest(&linkv1.GetPreviewRequest{})
			_, err := urlServer.GetPreview(ctx, req)
			Expect(err).To(HaveOccurred())
			Expect(connect.CodeOf(err)).To(Equal(connect.CodeUnimplemented))
		})
	})

	Context("GenerateQRCode", func() {
		It("should return unimplemented", func() {
			req := connect.NewRequest(&linkv1.GenerateQRCodeRequest{})
			_, err := urlServer.GenerateQRCode(ctx, req)
			Expect(err).To(HaveOccurred())
			Expect(connect.CodeOf(err)).To(Equal(connect.CodeUnimplemented))
		})
	})

})
