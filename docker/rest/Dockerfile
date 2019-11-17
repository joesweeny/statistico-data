# Step 1
FROM golang:1.12.9-alpine AS builder
RUN apk update && apk add git
ENV GO111MODULE=on
WORKDIR /app
COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo ./cmd/rest

# Step 2
FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /opt
COPY --from=builder /app/rest .
COPY ./api ./opt/api
EXPOSE 80
CMD ["./rest"]
