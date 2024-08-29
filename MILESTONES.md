# Milestones

## Milestone 1: Get the basic application running on Kubernetes
- [x] The server has a Dockerfile (server/Dockerfile)
- [x] The server runs on Kubernetes (kubernetes-manifests/linkserver.yaml)
- [x] The server is accessible from outside the cluster
- [x] The server is accessible from inside the cluster via port forwarding
- [x] Postgres is running on Kubernetes (kubernetes-manifests/postgres.yaml)
- [x] Postgres is accessible from outside the cluster
- [x] Postgres is accessible from inside the cluster via port forwarding
- [x] The server connects to Postgres
- [x] The server has an algorithm for shortening URLs which
- [x] The server can create a new short URL (curl /link.v1.UrlShortenerService/ShortenUrl) (server/cmd/main.go)
- [x] The server can retrieve a long URL from a short URL (curl /link.v1.UrlShortenerService/GetOriginalUrl) (server/cmd/main.go)
- [x] Running curl against the server returns the expected results (COMMANDS.md)

## Milestone 2: Performance Optimization
- [ ] Implement Redis caching layer for frequently accessed URLs
- [ ] Optimize database queries and add appropriate indexes
- [ ] Set up horizontal pod autoscaling for the linkservice deployment

## Milestone 3: High Availability and Scalability
- [ ] Implement a global load balancer for traffic distribution
- [ ] Configure database replication for improved read performance and failover
- [ ] Set up a multi-region Kubernetes cluster

## Milestone 4: User Authentication and Management
- [ ] Implement user registration and login functionality (create a new microservice for this, optionally use a new postgres instance)
- [ ] Implement API authentication (use redis to store tokens and validate them)
- [ ] Implement API rate limiting based on user authentication token 
- [ ] Add ability for users to create custom short URLs
- [ ] Implement a password reset functionality
- [ ] Implement URL expiration functionality
- [ ] Create user dashboard for managing shortened URLs

## Milestone 5: Web Frontend
- [ ] Set up a Svelte project using SvelteKit
- [ ] Integrate connectRPC for communication with the backend
- [ ] Implement responsive design for mobile and desktop
- [ ] Add user authentication flow in the frontend
- [ ] Create components for URL shortening and management

## Milestone 6: Analytics Pipeline
- [ ] Implement frontend API for tracking user interactions
- [ ] Set up Kafka for event streaming
- [ ] Develop a nearline stream processor for real-time analytics (use Flink)
- [ ] Implement data storage in MinIO
- [ ] Set up Apache Spark for batch processing of analytics data

## Milestone 7: Monitoring and Observability
- [ ] Implement distributed tracing with Zipkin
- [ ] Set up Prometheus for metrics collection and alerting
- [ ] Implement log aggregation with Loki
- [ ] Configure Grafana for unified visualization of metrics, logs, and traces
- [ ] Create custom dashboards and alerts for system health and performance

# Tech Stack

- Go
- Postgres
- Redis
- Kafka
- MinIO
- Apache Spark
- Svelte
- Kubernetes
- Zipkin
- Loki
- Grafana
- Prometheus
- Apache Spark