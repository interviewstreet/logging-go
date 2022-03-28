.DEFAULT_GOAL 	:= help

compile-app:
	@mkdir -p build
	@go build -o build/app app/logger.go

compile-core:
	@mkdir -p build
	@go build -o build/core core/logger.go

compile-request:
	@mkdir -p build
	@go build -o build/request request/logger.go

gomod_tidy:
	 go mod tidy

gofmt:
	go fmt -x ./...

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
