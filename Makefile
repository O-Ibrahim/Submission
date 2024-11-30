run:
	@go run cmd/main.go
clear-logs:
	rm *.log
build:
	@go build -o bin/RemoteCommands cmd/main.go

test:
	@export CGO_ENABLED=1
	@go test -v ./...
