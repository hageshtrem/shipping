
build:
	@protoc  --go_out=:. --go-grpc_out=:. proto/booking.proto

.PHONY: build