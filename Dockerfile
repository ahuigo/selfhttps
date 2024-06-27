FROM golang:1.22

MAINTAINER Ahuigo <github.com/ahuigo>

ENV GOPATH /go

COPY . /app/selfhttps
WORKDIR /app/selfhttps
RUN make install

ENTRYPOINT ["/go/bin/arun"]
