build:
	protoc --proto_path=. --go_out=plugins=grpc:. internal/app/proto/*.proto