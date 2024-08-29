# Stage 1: Generate protobuf code
FROM bufbuild/buf:1.34.0 as bufbuild
WORKDIR /app
COPY buf.gen.yaml ./
COPY proto ./proto
RUN buf generate proto

# Stage 2: Build the Go binary
FROM golang:1.22-alpine as builder
RUN apk add --no-cache git
WORKDIR /app

# Copy go.mod and go.sum files
COPY server/go.mod server/go.sum ./server/
COPY generated/go/go.mod generated/go/go.sum ./generated/go/

# Copy the generated protobuf code
COPY --from=bufbuild /app/generated ./generated

# Set WORKDIR to server directory
WORKDIR /app/server

# Download dependencies
RUN go mod download

# Add golang-migrate tool and migrations
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
COPY server/migrations ./migrations

# Copy the server source code
COPY server .

# Build the application
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go

# Final stage
FROM alpine:3.14

# Install ca-certificates for HTTPS requests and tzdata for timezone data
RUN apk add --no-cache ca-certificates tzdata

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/server/server .

# Copy the migration files
COPY --from=builder /app/server/migrations ./migrations

# Copy the migrate binary
COPY --from=builder /go/bin/migrate .

# Set the command to run the server
CMD ["./server"]