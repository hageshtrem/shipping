FROM rust:1.48 as builder
RUN rustup component add rustfmt
WORKDIR /app
COPY . /app
RUN cargo build --release --bin server

FROM debian:stable-slim
RUN apt-get update && \
    apt-get install -y openssl wget
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.2 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe
WORKDIR /app
COPY --from=builder /app/target/release/server ./handling-svc
CMD ["./handling-svc"]