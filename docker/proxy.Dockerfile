FROM golang:1.12 as builder

WORKDIR /root
COPY bin/proxy /bin/proxy
ENTRYPOINT [ "/bin/proxy" ]