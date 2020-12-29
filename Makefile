COVERAGE_FILE="coverage.out"

go.build:
	go build

GO_TEST=go test ./pkg/...

go.test.unit:
	$(GO_TEST) -race -covermode=atomic -coverprofile=$(COVERAGE_FILE)

go.test:
	$(GO_TEST) -race -covermode=atomic -func -coverprofile=$(COVERAGE_FILE)


GO_COVER=go tool cover -func=$(COVERAGE_FILE)

go.cover:
	$(GO_COVER)

go.cover.full: go.test go.cover

go.cover.total:
	$(GO_COVER) | grep total | awk '{print substr($$3, 1, length($$3)-1)}'

# make act ARGS="-s CODECOV_TOKEN=..."
act:
	docker build -t act-node-slim build/act/
	act -P ubuntu-latest=act-node-slim:latest $(ARGS)

install.lint:
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.33.0