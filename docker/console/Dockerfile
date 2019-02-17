# Step 1
FROM golang:latest AS builder
WORKDIR /go/src/github.com/joesweeny/statshub
COPY . .
RUN apt-get update && apt-get install git
RUN go version
RUN go get -u github.com/golang/dep/cmd/dep
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN dep ensure
RUN go install ./vendor/github.com/golang/protobuf/protoc-gen-go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo ./cmd/console

# Step 2
FROM alpine:3.8
RUN apk --no-cache add ca-certificates
WORKDIR /opt
COPY --from=builder /go/src/github.com/joesweeny/statshub/console .
CMD ["/bin/sh"]