# Milestones

## Milestone 1: The server runs on Kubernetes

- [x] The server has a Dockerfile (server/Dockerfile)
- [x] The server runs on Kubernetes (kubernetes-manifests/linkserver.yaml)
- [x] The server is accessible from outside the cluster
- [x] The server is accessible from inside the cluster via port forwarding

## Milestone 2: Postgres is running on Kubernetes (but not connected to the server)

- [x] Postgres is running on Kubernetes (kubernetes-manifests/postgres.yaml)
- [x] Postgres is accessible from outside the cluster
- [x] Postgres is accessible from inside the cluster via port forwarding

## Milestone 3: The Server is connected to Postgres

- [x] The server connects to Postgres
- [x] The server has an algorithm for shortening URLs which
- [x] The server can create a new short URL (curl /link.v1.UrlShortenerService/ShortenUrl) (server/cmd/main.go)
- [x] The server can retrieve a long URL from a short URL (curl /link.v1.UrlShortenerService/GetOriginalUrl) (server/cmd/main.go)
- [x] Running curl against the server returns the expected results (COMMANDS.md)