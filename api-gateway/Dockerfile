FROM golang:1.22-alpine as builder
RUN apk add --no-cache git
WORKDIR /app

# Copy go.mod and go.sum files
COPY api-gateway/go.mod api-gateway/go.sum ./api-gateway/
COPY generated/go/go.mod generated/go/go.sum ./generated/go/

# Copy the generated protobuf code
COPY generated ./generated

# Set WORKDIR to api-gateway directory
WORKDIR /app/api-gateway

# Download dependencies
RUN go mod download

# Copy the api-gateway source code
COPY api-gateway .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway ./cmd/main.go

FROM alpine:3.14
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/api-gateway/api-gateway .
CMD ["./api-gateway"]