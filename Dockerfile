FROM golang:1.18.6-alpine3.16 AS builder
WORKDIR /usr/app
COPY . .
ENV CGO_ENABLED=0 \
    GOOS=linux
RUN go get -d -t ./... && go test ./...
RUN go build -o server

FROM alpine:3.16.2
WORKDIR /usr/app
COPY --from=builder /usr/app/server /usr/app/schema.sql ./
CMD ["./server"]
