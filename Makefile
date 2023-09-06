.PHONY: run setup build vet test test-cov clean

run:
	@echo "Run the cli using 'go run' to pass arguments: 'go run ./cmd/jen <args>'"
	@echo "Ex: 'go run ./cmd/jen help'"

setup:
	@go mod download

build:
	@mkdir -p ~/bin
	@go build -ldflags "-X main.version=dev-$(date +%F-%T)" -o ~/bin/jen ./cmd/jen

vet:
	@go vet ./cmd/jen/...

test: vet
	@go test -count=1 ./cmd/jen/...

test-cov: vet
	@mkdir -p ./reports
	@go test ./cmd/jen/... -coverprofile ./reports/coverage.out -covermode count
	@go tool cover -func ./reports/coverage.out

clean:
	@rm -rf ./reports ./out