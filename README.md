# URL Shortener Service

This project implements a URL shortener service using Go, gRPC, and Kubernetes.

## Prerequisites

- [Go](https://golang.org/doc/install) (version 1.16 or later)
- [Docker](https://docs.docker.com/get-docker/)
- [Minikube](https://minikube.sigs.k8s.io/docs/start/)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Skaffold](https://skaffold.dev/docs/install/)
- [Buf](https://buf.build/docs/installation)

## Project Setup

1. Clone the repository:
   ```
   git clone https://github.com/wcygan/url-shortener.git
   cd url-shortener
   ```

2. Generate Proto files:
   ```
   buf generate proto
   ```

## Starting the Application

1. Start Minikube:
   ```
   minikube start
   ```

2. Use Skaffold to build and deploy the application:
   ```
   skaffold dev
   ```

   This command will build the Docker image, deploy it to Minikube, and set up port forwarding.

3. Wait for the deployment to complete. You should see output indicating that the service is running and port forwarding is set up.

## Testing the Service

Once the application is running, you can use curl commands to test the service:

### Shorten URL

Command:

```bash
curl \
    --header "Content-Type: application/json" \
    --data '{"original_url": "https://www.example.com/very/long/url"}' \
    http://localhost:8080/link.v1.UrlShortenerService/ShortenUrl
```

Response:

```bash
{"shortenedUrl":"http://short.url/6u8egU"}
```

### Get Original URL

Command:

```bash
curl \
    --header "Content-Type: application/json" \
    --data '{"shortened_url": "http://short.url/6u8egU"}' \
    http://localhost:8080/link.v1.UrlShortenerService/GetOriginalUrl
```

Response:

```bash
{"originalUrl":"https://www.example.com/very/long/url"}
```

Replace the example URLs with actual shortened URLs returned by the service.

## Stopping the Application

To stop the application and clean up resources:

1. Press `Ctrl+C` in the terminal where Skaffold is running.
2. Stop Minikube:
   ```
   minikube stop
   ```