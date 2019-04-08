FROM golang:1.11-alpine3.8

WORKDIR /go/src/github.com/statistico/statistico-data

RUN apk update && apk add git gcc g++
RUN go get gotest.tools/gotestsum

CMD ["/bin/sh"]