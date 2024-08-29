## Generate Proto

```bash
buf generate proto
```

## Start Server

From the root:

```bash
go run -C server cmd/main.go
```

From within the [server](server) directory:

```bash
go run cmd/main.go
```

## Curl

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