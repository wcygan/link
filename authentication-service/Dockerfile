FROM golang:1.22-alpine as builder
RUN apk add --no-cache git
WORKDIR /app

# Copy go.mod and go.sum files
COPY authentication-service/go.mod authentication-service/go.sum ./authentication-service/
COPY generated/go/go.mod generated/go/go.sum ./generated/go/

# Copy the generated protobuf code
COPY generated ./generated

# Set WORKDIR to authentication-service directory
WORKDIR /app/authentication-service

# Download dependencies
RUN go mod download

# Add golang-migrate tool and migrations
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
COPY authentication-service/migrations ./migrations

# Copy the authentication-service source code
COPY authentication-service .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/main.go

FROM alpine:3.14
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/authentication-service/auth-service .
COPY --from=builder /app/authentication-service/migrations ./migrations
COPY --from=builder /go/bin/migrate .
CMD ["./auth-service"]