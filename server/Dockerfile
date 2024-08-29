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

# Stage 3: Create the final minimal image
FROM alpine:3.20.1
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/server/server .
COPY --from=builder /app/server/migrations ./migrations
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
CMD ["./server"]