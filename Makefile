.PHONY: init
init:
	GO111MODULE=on go mod download

.PHONY: build
build:
	GO111MODULE=on go build

.PHONY: test
test:
	go test ./...

.PHONY: test-v
test-v:
	go test -v ./...

.PHONY: benchmark
benchmark:
	go test -bench . -benchmem

.PHONY: lint
lint:
	GO111MODULE=on golint ./...