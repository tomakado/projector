test:
	go test ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./... && \
	go tool cover -func=coverage.out

lint:
	golangci-lint run ./...
