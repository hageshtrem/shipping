FROM golang:1.15.0-alpine3.12 as builder
WORKDIR /go/src/pathfinder/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pathfinder ./cmd/server/main.go

FROM alpine:3.12
WORKDIR /app
COPY --from=builder /go/src/pathfinder/pathfinder .
ARG PORT
ENV PORT $PORT
EXPOSE $PORT
CMD ["./pathfinder"]