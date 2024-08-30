### Stage 1: Generate protobuf code ###
FROM bufbuild/buf:1.34.0 as bufbuild
WORKDIR /app
COPY buf.gen.yaml ./
COPY proto ./proto
RUN buf generate proto

### Stage 2: Build the Go binary ###
FROM golang:1.22-alpine as builder
RUN apk add --no-cache git
WORKDIR /app

# Copy go.mod and go.sum files
COPY link-service/go.mod link-service/go.sum ./link-service/
COPY generated/go/go.mod generated/go/go.sum ./generated/go/

# Copy the generated protobuf code
COPY --from=bufbuild /app/generated ./generated

# Set WORKDIR to link-service directory
WORKDIR /app/link-service

# Download dependencies
RUN go mod download

# Add golang-migrate tool and migrations
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
COPY link-service/migrations ./migrations

# Copy the link-service source code
COPY link-service .

# Build the application
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o link-service ./cmd/main.go

### Final stage ###
FROM alpine:3.14
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/link-service/link-service .
COPY --from=builder /app/link-service/migrations ./migrations
COPY --from=builder /go/bin/migrate .
CMD ["./link-service"]