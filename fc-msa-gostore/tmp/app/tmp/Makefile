run:
	go run $(name)/cmd/main.go

build:
	go build -o bin/$(name) $(name)/cmd/main.go

fmt:
	go fmt ./...

proto.pre:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

proto.gen:
	protoc --go_out=. --go-grpc_out=. pkg/proto/v1/$(name).proto