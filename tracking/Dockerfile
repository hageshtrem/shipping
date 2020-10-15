FROM golang:1.15.0-alpine3.12 as builder
WORKDIR /go/src/tracking/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tracking-svc ./cmd/server/tracking.go

FROM alpine:3.12
WORKDIR /app
COPY --from=builder /go/src/tracking/tracking-svc .
CMD ["./tracking-svc"]