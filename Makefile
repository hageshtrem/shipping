
gen: gen-booking gen-pathfinder

gen-pathfinder:
	@protoc  --go_out=:. --go-grpc_out=:. proto/pathfinder.proto

gen-booking:
	@protoc  --go_out=:. --go-grpc_out=:. proto/booking.proto


.PHONY: gen, gen-booking, gen-pathfinder