# Step 1
FROM golang:1.11-alpine3.8 AS builder
WORKDIR /go/src/github.com/joesweeny/statshub
COPY . .
RUN apk update && apk add git
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo ./cmd/grpc

# Step 2
FROM alpine:3.8
RUN apk --no-cache add ca-certificates
WORKDIR /opt
COPY --from=builder /go/src/github.com/joesweeny/statshub/grpc .
EXPOSE 50051
ENTRYPOINT ["./grpc", "--port 5001"]
CMD ["./grpc"]
