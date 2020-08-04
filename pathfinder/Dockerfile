FROM golang:latest as builder
WORKDIR /go/src/pathfinder/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o goapp ./cmd/server/main.go

FROM alpine:3.7
WORKDIR /app
COPY --from=builder /go/src/pathfinder/goapp .
ARG PORT
ENV PORT $PORT
EXPOSE $PORT
CMD ["./goapp"]