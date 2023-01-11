FROM golang:1.19.5-alpine3.16 AS builder
WORKDIR /usr/app
COPY . .
ENV CGO_ENABLED=0 \
    GOOS=linux
RUN go get ./... && go build -o server

FROM alpine:3.17
WORKDIR /usr/app
COPY --from=builder /usr/app/server /usr/app/schema.sql ./
EXPOSE 6001
CMD ["./server"]
