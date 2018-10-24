.SILENT:

unit_test:
	echo "\nunit tests\n"
	go test ./... -tags=unit 

acceptance_test:
	echo "\nacceptance tests\n"
	go test ./... -tags=acceptance 

test: unit_test acceptance_test

mocks:
	go generate ./...
