FROM golang:1.22

LABEL maintainer="Ahuigo <github.com/ahuigo>"

ENV GOPATH /go
RUN apt-get update && apt-get install -y openssl
COPY . /app/selfhttps
WORKDIR /app/selfhttps
RUN make install

ENTRYPOINT ["/go/bin/arun"]
