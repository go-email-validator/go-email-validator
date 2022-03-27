go.build:
	go build

GO_TEST=go test ./pkg/...
COVERAGE_UNIT_FILE="coverage.unit.out"
COVERAGE_FILE="coverage.out"
COVERAGE_TMP_FILE="coverage.out.tmp"
MOCK_PATTERN=mock_.*.go

GOROOT=$(shell go env GOROOT)

go.test.unit:
	$(GO_TEST) -race -covermode=atomic -coverprofile=$(COVERAGE_TMP_FILE)
	rm $(COVERAGE_UNIT_FILE) || true
	cat $(COVERAGE_TMP_FILE) | grep -v $(MOCK_PATTERN) > $(COVERAGE_UNIT_FILE)
	rm $(COVERAGE_TMP_FILE) || true

go.test:
	$(GO_TEST) -race -covermode=atomic -func -coverprofile=$(COVERAGE_TMP_FILE)
	rm $(COVERAGE_FILE) || true
	cat $(COVERAGE_TMP_FILE) | grep -v $(MOCK_PATTERN) > $(COVERAGE_FILE)
	rm $(COVERAGE_TMP_FILE) || true

go.mocks: go.mocks.isntall go.mocks.gen

go.mocks.isntall:
	go get github.com/golang/mock/mockgen

go.mocks.gen:
	mockgen -source=test/mock/net/interface.go -destination=test/mock/net/conn.go --package=mocknet
	mockgen -source=pkg/ev/evcache/evcache.go -destination=test/mock/ev/evcache/evcache.go --package=mockevcache
	mockgen -source=pkg/ev/evsmtp/smtpclient/interface.go -destination=test/mock/ev/evsmtp/smtpclient/interface.go --package=mocksmtpclient
	mockgen -source=pkg/ev/evsmtp/smtp.go -destination=pkg/ev/evsmtp/mock_smtp.go --package=evsmtp # https://github.com/golang/mock/issues/352
	mockgen -source=pkg/ev/validator.go -destination=test/mock/ev/validator.go --package=mockev

# wait https://github.com/golang/mock/issues/510
go.mocks.nets:
	mockgen -source=$(GOROOT)/go1.15.6/src/net/net.go -destination=test/mock/net/Conn.go --package=mocknet

go.lint:
	golint ./... | grep -v vendor/

GO_COVER=go tool cover -func=$(COVERAGE_FILE)
go.cover:
	$(GO_COVER)

go.cover.full: go.test go.cover

go.cover.total:
	$(GO_COVER) | grep total | awk '{print substr($$3, 1, length($$3)-1)}'

# make act ARGS="-s CODECOV_TOKEN=..."
act.build:
	docker build -t act-node-slim build/act/
act.run:
	act -P ubuntu-latest=act-node-slim:latest $(ARGS)

act: act.build act.run

install.lint:
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.33.0

mod.update: mod.update.main mod.update.as-email-verifier.adapter

mod.update.main:
	go mod tidy

mod.update.as-email-verifier.adapter:
	cd pkg/presentation/as-email-verifier/adapter && \
	go mod tidy

# Mount own self in go path
LIBRARY_PATH := $(GOROOT)/src$(GITHUB_VENDOR)

mount.for_adapter: umount.for_adapter
	rm -fr $(LIBRARY_PATH)
	mkdir -p $(LIBRARY_PATH)
	sudo mount -Br ~/go/src/github.com/go-email-validator/ $(LIBRARY_PATH)

umount.for_adapter:
	sudo umount $(LIBRARY_PATH) -q | exit 0


test.run.proxy:
	docker run --name proxy -dit --rm -p 1080:1080 -e 'SSS_USERNAME=username' -e 'SSS_PASSWORD=password' dijedodol/simple-socks5-server