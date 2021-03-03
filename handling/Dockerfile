FROM rust:1.48 as builder
RUN rustup component add rustfmt
WORKDIR /app
COPY . /app
RUN cargo build --release --bin server

FROM debian:stable-slim
RUN apt-get update && apt-get install -y openssl
WORKDIR /app
COPY --from=builder /app/target/release/server ./handling-svc
CMD ["./handling-svc"]