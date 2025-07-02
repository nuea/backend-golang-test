# Backend Golang Test

## Getting started

1. install main tool for develop
    ```shell
    brew install pre-commit
    brew install golangci-lint
    brew upgrade golangci-lint
    brew install buf
   
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    ```

2. pre-run before run project
    ```shell
    go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
    go get google.golang.org/protobuf/cmd/protoc-gen-go
    # go mod tidy - just install lib and clean unused lib in project
    make tidy
    # build && generate dependency injection
    make wire
    # generate everything eg. wire, proto-file
    make generate
    ```