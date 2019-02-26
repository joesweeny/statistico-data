FROM golang:1.11-alpine3.8

WORKDIR /go/src/github.com/joesweeny/statistico-data

RUN apk update && apk add git gcc g++
RUN go get gotest.tools/gotestsum

COPY ./ /go/src/github.com/joesweeny/statistico-data/

CMD ["/bin/sh"]