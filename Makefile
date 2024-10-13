# Detect the platform
ifeq ($(shell uname), Darwin)
    BROWSER_OPEN = open
else
    BROWSER_OPEN = xdg-open
endif

lint:
	gofumpt -w ./..
	golangci-lint run --fix

generate:
	go generate ./...

test:
	go test -v ./...

# Run all tests and generate coverage report in HTML
test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	$(BROWSER_OPEN) coverage.html
	go tool cover -func=coverage.out

client-run-local:
	POW_COMPLEXITY=4	CLIENT_HTTP_HOST_ADDR=:9091	SERVER_ADDR=http://localhost:9090	VALIDATE_ROUTE=/validate	CHALLENGE_ROUTE=/challenge	go run ./cmd/client/main.go   

server-run-local:
	SERVER_HTTP_HOST_ADDR=:9090	POW_COMPLEXITY=4	go run ./cmd/server/main.go   
					
run:
	docker-compose up --remove-orphans --build