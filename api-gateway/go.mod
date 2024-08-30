module github.com/wcygan/link/api-gateway

go 1.22.3

require (
	connectrpc.com/connect v1.16.2
	github.com/wcygan/link/generated/go v0.0.0
	golang.org/x/net v0.28.0
)

require (
	connectrpc.com/grpcreflect v1.2.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)

replace github.com/wcygan/link/generated/go => ../generated/go
