go.build:
	go build

go.test:
	go test ./... -race -covermode=atomic
