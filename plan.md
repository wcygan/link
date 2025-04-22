# PLAN.md - Link Service Development Plan

## 1. Overview

This document outlines the development plan for the Link Service project. The goal is to implement a service that shortens URLs, generates previews, and creates QR codes based on the requirements in `PRD.md` and using the technologies defined in `TECH-STACK.md`.

## 2. Current Status (as of analysis)

*   **API Definition:** `proto/link/v1/link.proto` defines the service contract (`UrlService` with RPCs: `ShortenURL`, `RefreshPreview`, `GetPreview`, `GenerateQRCode`). Generated Go code exists in `gen/`.
*   **Server Skeleton:** `server/cmd/main.go` sets up a basic ConnectRPC server, registers the `UrlService` handler, and enables gRPC reflection.
*   **Service Implementation:** `server/internal/service/url_service.go` contains the `UrlServer` struct, but all RPC methods are currently **unimplemented stubs**.
*   **Database:**
    *   TiDB is configured in `docker-compose.yaml`.
    *   SQL migrations for `urls` and `url_previews` tables exist in `server/db/migrations/`.
    *   **No database connection or interaction logic (GORM)** is implemented in the Go code yet.
*   **Caching:** Dragonfly is configured in `docker-compose.yaml`, but **not used** in the code.
*   **Workflow Engine:** Temporal is listed in `TECH-STACK.md` and detailed idempotency notes (`notes/idempotency.md`) rely on it, but the **Temporal SDK is not installed** and no workflows/activities are defined or used.
*   **Dependencies:** Core ConnectRPC and testing libraries (Ginkgo/Gomega) are present in `go.mod`. **Missing dependencies** include GORM, Temporal SDK, Viper, Zap/slog (if not using std `log`), UberFx, Goquery.
*   **Testing:** Basic test setup exists (`url_service_test.go`) confirming unimplemented status. **No actual logic tests exist**.
*   **Infrastructure:** `Makefile` provides useful commands for building, testing, proto generation, and Docker management.
*   **Configuration, Logging, DI:** Basic std `log` used. No configuration management (Viper) or Dependency Injection (UberFx) implemented yet.
*   **Design Notes:** `notes/idempotency.md` provides a valuable, detailed plan for handling idempotency using DB constraints, Temporal, and API patterns. `LAYOUT-AND-BEST-PRACTICES.md` provides guidelines.

## 3. Development Phases

### Phase 0: Foundation & Setup

**Goal:** Prepare the Go application structure for implementing features.

*   [x] **Install Dependencies:** Add missing Go modules specified in `TECH-STACK.md`.
    *   [x] GORM (`gorm.io/gorm` and TiDB driver `gorm.io/driver/mysql`)
    *   [x] Temporal SDK (`go.temporal.io/sdk`)
    *   [x] Viper (`github.com/spf13/viper`) for configuration
    *   [x] Zap (`go.uber.org/zap`) or configure standard `log/slog` for structured logging
    *   [x] UberFx (`go.uber.org/fx`) for dependency injection
    *   [x] Goquery (`github.com/PuerkitoBio/goquery`) for HTML parsing (needed for Previews)
    *   [x] QR Code library (e.g., `github.com/skip2/go-qrcode`)
    *   [x] Run `go mod tidy`
*   [ ] **Configuration:**
    *   Set up Viper (`server/internal/config/`) to load settings (e.g., DB DSN, server address, Temporal details) from environment variables or a config file.
*   [ ] **Logging:**
    *   Initialize structured logging (Zap or slog) early in `main.go`.
    *   Prepare to inject logger instances where needed.
*   [ ] **Dependency Injection (UberFx):**
    *   Refactor `main.go` to use UberFx.
    *   Define providers for configuration, logger, database connection, Temporal client, and the `UrlServer`.
*   [ ] **Database Connection:**
    *   Implement DB connection logic using GORM and the TiDB driver, managed by UberFx.
    *   Inject the `*gorm.DB` instance into the `UrlServer` (or a dedicated repository/store layer).
    *   Ensure migrations can be run (e.g., add a `make migrate-up` target or integrate auto-migration cautiously).

### Phase 1: Core URL Shortening (Database Interaction)

**Goal:** Implement the `ShortenURL` RPC endpoint, focusing on database persistence and idempotency. Refer heavily to `notes/idempotency.md`.

*   [ ] **Implement `UrlServer.ShortenURL`:**
    *   [ ] Receive `original_url`, `custom_code`, `ttl_seconds` from the request.
    *   [ ] **Canonicalize URL:** Implement or use a library for URL normalization (as discussed in `idempotency.md`).
    *   [ ] **Generate Short Code:** Implement logic to generate unique short codes if `custom_code` is not provided.
    *   [ ] **Database Interaction (GORM):**
        *   Check if the canonical URL already exists (`UNIQUE` constraint on `canonical_url`). If yes, return the existing short code.
        *   If `custom_code` is provided, check if it exists (`PRIMARY KEY` or `UNIQUE` constraint on `code`). Handle conflicts appropriately (e.g., return error).
        *   Attempt to `INSERT` the new record (`code`, `canonical_url`, `expires_at` based on `ttl_seconds`). Handle potential race conditions/duplicate errors gracefully (e.g., re-query if insert fails on unique constraint violated).
        *   Use `INSERT IGNORE` or `INSERT...ON DUPLICATE KEY UPDATE` strategy as outlined in `idempotency.md` for robustness.
    *   [ ] **Construct Response:** Return the `short_url` (requires knowing the base domain, fetch from config) and `short_code`.
*   [ ] **Idempotency (API Layer):**
    *   [ ] (Optional but recommended) Implement handling of `Idempotency-Key` header using Dragonfly cache, as suggested in `notes/idempotency.md`.
*   [ ] **Testing:**
    *   [ ] Write unit tests for URL canonicalization and code generation logic.
    *   [ ] Write **integration tests** (`url_service_test.go`) for `ShortenURL` using `testcontainers-go` to spin up a test TiDB instance. Test:
        *   Creating a new short URL.
        *   Retrieving an existing URL.
        *   Using a custom code.
        *   Handling conflicts (duplicate URL, duplicate custom code).
        *   TTL/Expiration logic (if feasible to test).
        *   Idempotency scenarios.

### Phase 2: URL Preview (Async Workflow with Temporal)

**Goal:** Implement background URL preview generation using Temporal and the corresponding API endpoints.

*   [ ] **Define Temporal Components:**
    *   [ ] **Activities:**
        *   `FetchURLMetadata`: Takes a URL, fetches HTML using Go's `net/http`, parses metadata (title, description, image) using `goquery`.
        *   `PersistPreview`: Takes preview data and `url_code`, saves it to the `url_previews` table using GORM.
    *   [ ] **Workflows:**
        *   `PreviewWorkflow`: Orchestrates fetching and persisting. Takes `url_code` and `canonical_url` as input. Calls `FetchURLMetadata` then `PersistPreview`. Use appropriate retry policies. Workflow ID: `"preview::" + url_code`.
        *   `(Optional)` `ShortLinkWorkflow`: As suggested in `idempotency.md`, could be triggered by `ShortenURL`, persist the URL row itself (though maybe redundant if done in handler), and then start `PreviewWorkflow` as a child workflow. Workflow ID: `"shortlink::" + canonical_url`.
*   [ ] **Setup Temporal Client & Worker:**
    *   Initialize Temporal client via UberFx.
    *   Setup a Temporal worker in a separate process or goroutine in `main.go` (or a dedicated worker service) to host the workflows and activities.
*   [ ] **Integrate Temporal into `UrlServer`:**
    *   [ ] Modify `ShortenURL` to **asynchronously trigger** the `PreviewWorkflow` (or `ShortLinkWorkflow`) after successfully creating/retrieving the URL record. Use `ExecuteWorkflow` with appropriate options (`WorkflowIDReusePolicy`, etc. from `idempotency.md`).
    *   [ ] Implement `UrlServer.RefreshPreview`:**
        *   Accept `short_code` or `original_url`. Resolve to `url_code` and `canonical_url` via DB lookup.
        *   Trigger the `PreviewWorkflow` using `SignalWithStartWorkflow` to either start it or signal an existing run.
    *   [ ] Implement `UrlServer.GetPreview`:**
        *   Accept `short_code` or `original_url`. Resolve to `url_code` via DB lookup.
        *   Query the `url_previews` table using GORM based on `url_code`.
        *   Return the found preview data or `NotFound` error.
*   [ ] **Testing:**
    *   [ ] Unit tests for Activities (`FetchURLMetadata` parsing logic).
    *   [ ] Integration tests using Temporal's test framework (`testsuite`) for the Workflows.
    *   [ ] Update `url_service_test.go` integration tests (with testcontainers-go for DB) to verify `RefreshPreview` triggers workflows (mock Temporal client or use test framework) and `GetPreview` retrieves data.

### Phase 3: QR Code Generation

**Goal:** Implement the `GenerateQRCode` RPC endpoint.

*   [ ] **Implement `UrlServer.GenerateQRCode`:**
    *   [ ] Receive `short_code` and styling parameters (`size`, colors).
    *   [ ] **Database Lookup:** Query the `urls` table to verify the `short_code` exists (optional, but good practice). Get the `canonical_url` or construct the full short URL.
    *   [ ] Use a QR code library (e.g., `skip2/go-qrcode`) to generate the QR code PNG data based on the full short URL or canonical URL. Apply styling parameters.
    *   [ ] Return the `png_data` and `content_type` ("image/png").
*   [ ] **Testing:**
    *   [ ] Write unit/integration tests for `GenerateQRCode` to ensure codes are generated and parameters are applied. Verify output is valid PNG (can check magic bytes).

### Phase 4: Refinement & Polish

**Goal:** Improve robustness, performance, and observability.

*   [ ] **Caching (Dragonfly):**
    *   Implement caching for resolved `GetPreview` results.
    *   Implement/verify `Idempotency-Key` caching for `ShortenURL`.
    *   Connect to Dragonfly via UberFx.
*   [ ] **Error Handling:**
    *   Review error handling across all layers (API, service, DB, Temporal). Ensure proper ConnectRPC error codes are returned (`connect.NewError`).
    *   Handle Temporal errors gracefully.
*   [ ] **Observability:**
    *   Add basic metrics (e.g., request counts, durations) using Prometheus client library.
    *   Ensure structured logs provide sufficient detail for debugging.
    *   Consider adding tracing if needed.
*   [ ] **Security:**
    *   Review input validation (e.g., URL formats, size limits).
    *   Consider rate limiting if this were a public service.
*   [ ] **Documentation:**
    *   Update `README.md` with detailed usage instructions for all endpoints, including curl examples (or using `grpcurl`).
    *   Ensure comments in code are clear.

## 4. Next Steps (Immediate Focus)

1.  **Phase 0:** Complete all tasks in the Foundation & Setup phase. This is blocking for subsequent feature implementation. Prioritize installing dependencies, setting up DI, config, logging, and basic DB connection.
2.  **Phase 1:** Start implementing `ShortenURL` focusing on the core DB logic and idempotency based on `notes/idempotency.md`. Write the corresponding tests.

## 5. Ongoing Tasks

*   **Testing:** Continuously write and run unit and integration tests (`make test`). Maintain high coverage.
*   **Linting & Formatting:** Regularly run `make lint` and `make fmt`.
*   **Code Reviews:** Conduct code reviews for implemented features.
*   **Documentation:** Keep `README.md`, `PLAN.md`, and code comments up-to-date.
*   **Dependency Updates:** Periodically check for and update dependencies.
