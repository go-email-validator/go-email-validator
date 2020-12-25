go.build:
	go build

go.test.unit:
	go test ./... -race -covermode=atomic

go.test:
	go test ./... -race -covermode=atomic -func
