FROM golang:1.15.0-alpine3.12 as builder
WORKDIR /go/src/tracking/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tracking-svc ./cmd/server/tracking_server.go

FROM alpine:3.12
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.2 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe
WORKDIR /app
COPY --from=builder /go/src/tracking/tracking-svc .
CMD ["./tracking-svc"]