.SILENT:

test:
	go test ./... -tags=unit 

mocks:
	go generate ./...
