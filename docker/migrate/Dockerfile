# Step 1
FROM golang:1.11-alpine3.8 AS builder
WORKDIR /go/src/app
RUN apk update && apk add git gcc g++
RUN go get -u github.com/pressly/goose/cmd/goose

# Step 2
FROM alpine:3.8
RUN apk --no-cache add ca-certificates
WORKDIR /opt
COPY --from=builder /go/bin/goose /usr/local/bin
COPY ./bin ./bin
COPY ./database ./database
COPY ./docker/migrate/wait-for-it.sh .
CMD ["/bin/sh"]