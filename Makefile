all: vet test testrace

build:
	go build google.golang.org/grpc/...

clean:
	go clean -i google.golang.org/grpc/...

deps:
	GO111MODULE=on go get -d -v google.golang.org/grpc/...

proto:
	PATH="${PATH}:${HOME}/go/bin" protoc -I ./proto \
    --go_out ./proto --go_opt paths=source_relative \
    --go-grpc_out ./proto --go-grpc_opt paths=source_relative \
    ./proto/helloworld/helloworld.proto

test:
	go test -cpu 1,4 -timeout 7m google.golang.org/grpc/...

testsubmodule:
	cd security/advancedtls && go test -cpu 1,4 -timeout 7m google.golang.org/grpc/security/advancedtls/...
	cd security/authorization && go test -cpu 1,4 -timeout 7m google.golang.org/grpc/security/authorization/...

testrace:
	go test -race -cpu 1,4 -timeout 7m google.golang.org/grpc/...

testdeps:
	GO111MODULE=on go get -d -v -t google.golang.org/grpc/...

vet: vetdeps
	./vet.sh

vetdeps:
	./vet.sh -install

.PHONY: \
	all \
	build \
	clean \
	proto \
	test \
	testrace \
	vet \
	vetdeps