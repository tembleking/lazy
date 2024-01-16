
fmt:
    go run mvdan.cc/gofumpt@latest -w -l .

# Lints the project
lint:
    go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run --timeout 1h

# Tests if the project builds correctly
test-build:
    go build ./...
    go test -run ^$ ./... 1>/dev/null

# Runs all tests
test:
    go run github.com/onsi/ginkgo/v2/ginkgo -r -p --race --cover

