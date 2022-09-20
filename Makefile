proto-gen:
	protoc --proto_path=proto proto/*.proto --go_out=pb --go-grpc_out=pb

stubs-gen:
	protoc --proto_path=proto proto/*.proto --plugin=$(go env GOPATH)/bin/protoc-gen-go-grpc --go-grpc_out=. --go_out=.

server:
	go run cmd/server/server.go

client:
	go run cmd/client/client.go

client-evans:
	evans -r repl --host localhost --port 50051