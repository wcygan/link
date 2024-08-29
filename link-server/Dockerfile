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
COPY link-server/go.mod link-server/go.sum ./link-server/
COPY generated/go/go.mod generated/go/go.sum ./generated/go/

# Copy the generated protobuf code
COPY --from=bufbuild /app/generated ./generated

# Set WORKDIR to link-server directory
WORKDIR /app/link-server

# Download dependencies
RUN go mod download

# Add golang-migrate tool and migrations
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
COPY link-server/migrations ./migrations

# Copy the link-server source code
COPY link-server .

# Build the application
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o link-server ./cmd/main.go

### Final stage ###
FROM alpine:3.14
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/link-server/link-server .
COPY --from=builder /app/link-server/migrations ./migrations
COPY --from=builder /go/bin/migrate .
CMD ["./link-server"]