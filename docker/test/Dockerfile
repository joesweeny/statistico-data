FROM golang:1.11-alpine3.8
WORKDIR /go/src/github.com/joesweeny/statshub
RUN apk update && apk add git gcc g++
RUN go get -u github.com/satori/go.uuid
RUN go get -u github.com/lib/pq
CMD ["/bin/sh"]