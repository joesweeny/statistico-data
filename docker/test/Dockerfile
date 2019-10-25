FROM golang:1.12.9-alpine

RUN apk update && apk add git gcc g++
ENV GO111MODULE=on
WORKDIR /opt
COPY go.mod .
COPY go.sum .

CMD ["/bin/sh"]