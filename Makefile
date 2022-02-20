test:
	go test ./...

test-coverage:
	go test -v -covermode=count -coverprofile=coverage_raw.out ./... && \
	cat coverage_raw.out | grep -v pkg/create.go > coverage.out && \
	go tool cover -func=coverage.out

lint:
	golangci-lint run ./...
