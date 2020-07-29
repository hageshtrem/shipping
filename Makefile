
gen: gen-booking gen-pathfinder

gen-pathfinder:
	@protoc  --go_out=:. --go-grpc_out=:. proto/pathfinder.proto

gen-booking:
	@protoc  --go_out=:. --go-grpc_out=:. proto/booking.proto
	@protoc  --go_out=booking/pb --go-grpc_out=booking/pb proto/pathfinder.proto


.PHONY: gen, gen-booking, gen-pathfinder