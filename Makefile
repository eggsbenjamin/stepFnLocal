.SILENT:

test:
	go test ./... -tags=unit 

gen_mocks:
	go generate ./...
