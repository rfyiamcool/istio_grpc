DEFAULT_GOAL := build-all

build-all: clean gen build

gen:
	@ protoc -I/usr/local/include -I. -I$(GOPATH)/src \
	--go_out=plugins=grpc:. proto/call.proto

build:
	@go build -o bin/client cmd/client/client.go
	@go build -o bin/proxy cmd/proxy/proxy.go
	@go build -o bin/cache cmd/cache/cache.go
	@go build -o bin/store cmd/store/store.go

test:
	@go test -v 

clean:
	@rm -rf proto/call.pb.go
	@rm -rf bin/proxy
	@rm -rf bin/cache
	@rm -rf bin/store

run:
	# install foreman or goreman
	@goreman start

deps:
	@dep ensure

docker-build:
	@docker build -f docker/proxy.Dockerfile -t rfyiamcool/proxy .
	@docker build -f docker/cache.Dockerfile -t rfyiamcool/cache .
	@docker build -f docker/store.Dockerfile -t rfyiamcool/store .