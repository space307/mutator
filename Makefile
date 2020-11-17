.DEFAULT_GOAL := build

.PHONY: build
build:
	go mod verify
	go build

.PHONY: test-unit
test-unit:
	go clean -testcache
	go test -race -v -run Unit ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: lint-fix
lint-fix:
	goimports -local github.com/space307/mutator -w .
	go fmt ./...
	golangci-lint -v run

.PHONY: generate
generate:
	go generate ./...
