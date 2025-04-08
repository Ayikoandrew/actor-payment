run: build
	@./bin/ap

build:
	@go build -o ./bin/ap

test:
	@go test -v ./...

proto:
	@protoc --go_out=./msg --go_opt=paths=source_relative \
	--go-grpc_out=./msg --go-grpc_opt=paths=source_relative \
	msg/messages.proto