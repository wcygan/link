# Milestones

## Milestone 1: The server runs on Kubernetes

- [x] The server has a Dockerfile (server/Dockerfile)
- [x] The server runs on Kubernetes (kubernetes-manifests/linkserver.yaml)
- [x] The server is accessible from outside the cluster
- [x] The server is accessible from inside the cluster via port forwarding

## Milestone 2: Postgres is running on Kubernetes (but not connected to the server)

- [ ] Postgres is running on Kubernetes (kubernetes-manifests/postgres.yaml)
- [ ] Postgres is accessible from outside the cluster
- [ ] Postgres is accessible from inside the cluster via port forwarding

## Milestone 3: The Server is connected to Postgres

- [ ] The server connects to Postgres
- [ ] The server has an algorithm for shortening URLs which
- [ ] The server can create a new short URL (curl /link.v1.UrlShortenerService/ShortenUrl) (server/cmd/main.go)
- [ ] The server can retrieve a long URL from a short URL (curl /link.v1.UrlShortenerService/GetOriginalUrl) (server/cmd/main.go)
- [ ] Running curl against the server returns the expected results (COMMANDS.md)