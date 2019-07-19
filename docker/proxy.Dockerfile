FROM golang:1.12 as builder

RUN apk add --no-cache make

WORKDIR /root
ADD . /go/src/github.com/rfyiamcool/istio_grpc
RUN cd /go/src/github.com/rfyiamcool/istio_grpc && make build && mv bin/* /bin

ENTRYPOINT [ "/bin/proxy" ]