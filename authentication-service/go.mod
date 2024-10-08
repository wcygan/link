module github.com/wcygan/link/authentication-service

go 1.22.3

toolchain go1.22.6

require (
	connectrpc.com/connect v1.16.2
	github.com/jackc/pgx/v5 v5.5.5
	github.com/wcygan/link/generated/go v0.0.0
	golang.org/x/crypto v0.26.0
	golang.org/x/net v0.28.0
)

require (
	connectrpc.com/grpcreflect v1.2.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)

replace github.com/wcygan/link/generated/go => ../generated/go
