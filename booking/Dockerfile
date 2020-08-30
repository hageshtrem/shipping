FROM golang:1.15.0-alpine3.12 as builder
WORKDIR /go/src/booking/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o booking-svc ./cmd/server/booking.go

FROM alpine:3.12
WORKDIR /app
COPY --from=builder /go/src/booking/booking-svc .
CMD ["./booking-svc"]