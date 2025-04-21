# Link Service

Shorten, preview, and generate QR codes for URLs.

## Running the Server

To run the Go server locally:

```bash
go run ./server/cmd/main.go
```

The server will start on `localhost:8080`.

### Using gRPC UI

Once the server is running, you can interact with it using `grpcui`. If you don't have it, install it:

```bash
go install github.com/fullstorydev/grpcui/cmd/grpcui@latest
```

Then, run `grpcui` pointing to the server:

```bash
grpcui -plaintext localhost:8080
```

This will open a web interface in your browser to interact with the API.