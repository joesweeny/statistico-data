# Step 1
FROM golang:1.11-alpine3.8 AS builder
WORKDIR /go/src/github.com/joesweeny/statistico-data/cmd/api
COPY . .
RUN apk update && apk add git
RUN go get github.com/gorilla/mux
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo ./cmd/api

# Step 2
FROM alpine:3.8
RUN apk --no-cache add ca-certificates
WORKDIR /opt
COPY --from=builder /go/src/github.com/joesweeny/statistico-data/cmd/api/api .
EXPOSE 8000
ENTRYPOINT ["./api", "--port 8000"]
CMD ["./api"]
