run: gen
	docker-compose up -d

gen: gen-booking gen-pathfinder gen-tracking

gen-pathfinder:
	@protoc  --go_out=:. --go-grpc_out=:. proto/pathfinder.proto

gen-booking:
	@protoc  --go_out=:. --go-grpc_out=:. proto/booking.proto
	@protoc  --go_out=:. proto/booking_events.proto
	@protoc  --go_out=booking/pb --go-grpc_out=booking/pb proto/pathfinder.proto

gen-tracking:
	@protoc  --go_out=:. --go-grpc_out=:. proto/tracking.proto
	@protoc  --go_out=tracking/pb proto/booking_events.proto

.PHONY: gen gen-booking gen-pathfinder gen-tracking run