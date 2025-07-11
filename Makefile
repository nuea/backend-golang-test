HELP_CMD = grep -E '^[a-zA-Z_-]+:.*?\#\# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?\#\# "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: tidy generate ## build for use on image
	@cd cmd/http && go build -o ../../.bin/http .
	@cd cmd/grpc && go build -o ../../.bin/grpc .

tidy: ## download dependencies 
	go mod tidy

wire: ## genarate google wire
	go run github.com/google/wire/cmd/wire ./...

generate: ## go generate
	go mod tidy
	go generate ./...

lint: ## lint
	golangci-lint run

http: ## run http server
	GIN_MODE=debug ENV=local go run cmd/http/main.go

grpc: ## run grpc server
	GO111MODULE=on ENV=local go run cmd/grpc/main.go

help: ## show this help
	@${HELP_CMD}

proto-libs: ## proto - install libs last version
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/golang/mock/mockgen@latest
	go mod tidy

test:
	go test -coverprofile=c.out ./...
	go tool cover -html=c.out -o=coverage.html
