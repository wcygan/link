# Guaranteeing Idempotency in URL Shortener

## Goals Recap

*   One-to-one mapping between each canonical "real" URL and a short code.
*   API must tolerate retries (network failure, client refresh) without making duplicates.
*   Callers may supply either the original URL or the short code in follow-up requests.

Below is a prescriptive plan that enforces those rules at every layer—database, Temporal, and ConnectRPC.

## 1 · Canonicalise first—then guarantee uniqueness

| Step | What to do                                                                                                                                                           | Why it matters                                                                   |
| :--- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------- |
| a.   | Canonicalise the URL                                                                                                                                                 | Lower-case scheme+host, remove default port, strip trailing /, sort query params, etc. Adopt an RFC-compliant normaliser (e.g., github.com/peterbourgon/uritemplates). | Two superficially different strings (HTTPS://ExaMple.com, https://example.com/) map to one row. |
| b.   | Store `canonical_url`                                                                                                                                              | Keep a separate `canonical_url` column and add a `UNIQUE` index on it.           | The DB—not your code—enforces 1:1 mapping.                                       |
| c.   | Short code uniqueness                                                                                                                                                  | Primary Key or another `UNIQUE` index on `code`.                                 | Protects custom codes and auto-generated ones.                                   |

```sql
CREATE TABLE urls (
  code           VARCHAR(32)  PRIMARY KEY,
  canonical_url  VARCHAR(2048) UNIQUE NOT NULL,
  created_at     TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
  expires_at     TIMESTAMP    NULL,
  TTL = 'expires_at + INTERVAL 0 SECOND'
);
```

*TiDB guarantees uniqueness across the cluster; optimistic transactions automatically retry on write-write conflicts.*

## 2 · Idempotent ShortLinkWorkflow

| Temporal primitive      | Setting                                                                      | Effect                                                                                                    |
| :---------------------- | :--------------------------------------------------------------------------- | :-------------------------------------------------------------------------------------------------------- |
| Workflow ID             | `shortlink::<canonical_url>`                                                 | Starting the same link twice hits the `WorkflowExecutionAlreadyStarted` error instead of duplicating work. |
| `WorkflowIdReusePolicy` | `AllowDuplicateFailedOnly`                                                   | Retries succeed if a previous run failed, but no second run starts while one is still open/succeeded.     |
| `PersistUrl` Activity   | `INSERT IGNORE INTO urls ...` or `INSERT … ON DUPLICATE KEY UPDATE code=code` | First writer wins; on conflict we `SELECT code` and return it. TiDB's unique index ensures safety.        |

*Retry loop: if a race causes duplicate-code error (rare when you auto-generate codes), the Activity just generates a new code and retries—because the canonical URL key is already locked, no second row is created.*

*The parent workflow always starts (or re-attaches to) `PreviewWorkflow` using `preview::<code>` as its child ID.*

## 3 · Idempotent PreviewWorkflow

*   **Workflow ID:** `preview::<code>`

*If called through `RefreshPreview` RPC, we attempt `workflow.Start` with the same ID—Temporal turns that into a Signal (`StartSignalWithStart`) if the Workflow is already running, so you refresh immediately without duplicate histories.*

## 4 · API design for flexible inputs

### 4.1 Proto

```protobuf
message UrlSelector {
  oneof by {
    string short_code = 1;
    string url        = 2;
  }
}

service UrlService {
  rpc ShortenURL     (ShortenURLRequest) returns (ShortenURLResponse);
  rpc RefreshPreview (UrlSelector)       returns (google.protobuf.Empty);
  rpc GetPreview     (UrlSelector)       returns (UrlPreview);
  rpc GenerateQRCode (UrlSelector)       returns (QrCode);
}
```

*`oneof` forces callers to supply exactly one identifier—no ambiguity.*

*Handlers canonicalise `url` (if provided) ➜ lookup row ➜ derive `code`, or vice-versa.*

### 4.2 HTTP + ConnectRPC idempotency

*   Encourage clients to send an `Idempotency-Key` header (a UUID) on `ShortenURL` requests.
*   Store the key in a short-lived in-memory cache (e.g., Dragonfly) mapping ➜ `canonical_url`.
*   If the same key re-appears, immediately return the cached response.

*This makes accidental double-posts (browser retry, mobile reconnect) resolve instantly, while the database/Temporal layers remain the ultimate safety net.*

## 5 · End-to-end flow (happy path)

1.  `POST /ShortenURL` with `Idempotency-Key: abc123`, `URL=https://Example.com/`.
2.  Handler canonicalises → `https://example.com`.
3.  Attempts `INSERT`.
4.  If row exists → returns existing `code`.
5.  Else inserts new row (`code=KfX2a`) and starts `ShortLinkWorkflow`.
6.  Parent workflow persists row (already done), returns `KfX2a`, starts `PreviewWorkflow`.
7.  If the client retries with the same header or the same canonical URL, step 3b short-circuits and returns the same result—no duplicate row, no duplicate Workflow.

## Key Take-aways

*   DB unique indexes + canonicalisation = bullet-proof 1:1 mapping.
*   Temporal Workflow IDs lock each logical operation to a single execution.
*   `Idempotency-Key` header gives instant deduplication at the API edge for user experience.
*   `oneof` selector in Protobuf lets every downstream call accept either identifier without ambiguity.

*Follow these rules and you get strict referential integrity, race-free idempotency, and a developer-friendly API—all while staying completely stateless in your Go handlers.*

