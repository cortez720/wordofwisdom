
lint:
	gofumpt -w ./..
	golangci-lint run --fix

generate:
	go generate ./...

test:
	go test -v ./...

client-run:
	go run ./cmd/client/main.go   

server-run:
	go run ./cmd/server/main.go   