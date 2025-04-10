run: build
	@./bin/broker & ./bin/processor

build:
	@mkdir -p ./bin ./tmp
	@go build -o ./bin/broker ./broker/main.go
	@go build -o ./bin/processor ./cmd/processor/main.go
	@cp ./bin/broker ./tmp/broker
	@cp ./bin/processor ./tmp/processor

dev:
	@docker compose down && docker compose up --build

test:
	@go test -v ./...

proto:
	@protoc --go_out=./msg --go_opt=paths=source_relative \
	--go-grpc_out=./msg --go-grpc_opt=paths=source_relative \
	msg/messages.proto