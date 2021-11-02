proto:
	protoc --proto_path=./proto --go_out=:. --go-grpc_out=:. ./proto/user.proto	
clean-all:  
	go clean -i ./...
fmt:    
	go fmt ./...
	goimports -w $(FILES)
test:
	go test -v ./... -short
start:  
	docker compose up
