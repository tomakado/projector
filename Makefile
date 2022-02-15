test:
	go test ./...

test-coverage:
	go test -coverprofile=coverage.out ./... && \
	go tool cover -func=coverage.out
