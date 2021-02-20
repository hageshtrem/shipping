FROM rust:1.48 as builder
RUN rustup component add rustfmt
RUN cargo install wasm-pack
RUN cargo install cargo-make
WORKDIR /app
COPY . /app
# Needs to set env API_URL="http://localhost/v1"
# ARG API_URL
# ENV API_URL $API_URL
RUN cargo make build_release

FROM nginx:1.19.4
COPY --from=builder /app/pkg /usr/share/nginx/html/pkg
COPY --from=builder /app/index.html /usr/share/nginx/html
