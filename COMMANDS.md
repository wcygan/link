## Generate Proto

```bash
buf generate proto
```

## Start Server

Run within the [server](server) directory:

```bash
go run server/cmd/main.go
```

## Curl

### Shorten URL

```
curl \
    --header "Content-Type: application/json" \
    --data '{"original_url": "https://www.example.com/very/long/url"}' \
    http://localhost:8080/link.v1.UrlShortenerService/ShortenUrl
```

### Get Original URL

```
curl \
    --header "Content-Type: application/json" \
    --data '{"shortened_url": "http://short.url/abcde"}' \
    http://localhost:8080/link.v1.UrlShortenerService/GetOriginalUrl
```