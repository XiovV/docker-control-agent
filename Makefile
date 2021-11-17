.PHONY: audit
audit:
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	golangci-lint run
	@echo 'Running tests...'
	go test ./...

.PHONY: codecov
codecov:
	@echo 'Generating code coverage report...'
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: protoc
protoc:
	@echo 'Compiling .proto'
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpc/dokkup.proto