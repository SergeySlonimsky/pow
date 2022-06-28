build-client:
	go build -o app ./cmd/client

build-server:
	go build -o app ./cmd/server

format:
	find . -type f -name '*.go' -not -path "./vendor/*" -not -path "*_mock.go" \
		-exec goimports -local "github.com/SergeySlonimsky/pow" -w {} \; \
		-exec gci write -s Std -s Default -s "Prefix(github.com/SergeySlonimsky/pow)" {} \;
	gofumpt -l -w .

lint:
	golangci-lint run --allow-parallel-runners

generate-go:
	go generate ./...

test:
	go test ./...
