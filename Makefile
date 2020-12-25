run: gen
	docker-compose up -d

gen: gen-booking gen-pathfinder gen-tracking gen-handling

gen-pathfinder:
	@protoc --go_out=:. --go-grpc_out=:. proto/pathfinder.proto

gen-booking:
	@protoc --proto_path=proto --go_out=:. --go-grpc_out=:. proto/booking.proto proto/itinerary.proto
	@protoc --proto_path=proto --go_out=:. proto/booking_events.proto proto/itinerary.proto
	@protoc --go_out=booking/pb --go-grpc_out=booking/pb proto/pathfinder.proto

gen-tracking:
	@protoc --go_out=:. --go-grpc_out=:. proto/tracking.proto
	@protoc --proto_path=proto --go_out=tracking/pb proto/booking_events.proto proto/itinerary.proto

gen-handling:
	@cp -r proto/ handling/proto

.PHONY: gen gen-booking gen-pathfinder gen-tracking gen-handling run
