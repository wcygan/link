module github.com/wcygan/link/server

go 1.22.3

require (
	connectrpc.com/connect v1.16.2
	github.com/wcygan/link/generated/go v0.0.0
	golang.org/x/net v0.28.0
)

require (
	github.com/bufbuild/connect-go v1.10.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)

replace github.com/wcygan/link/generated/go => ../generated/go
