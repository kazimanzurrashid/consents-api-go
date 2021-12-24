FROM golang:1.17.5-alpine3.15 AS builder
WORKDIR /usr/app
COPY . .
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go get -d -t ./... && go test ./...
RUN go build ./cmd/server && go build ./cmd/health-check

FROM scratch
WORKDIR /usr/app
COPY --from=builder /usr/app/server /usr/app/health-check /usr/app/schema.sql ./
EXPOSE 8080
HEALTHCHECK --start-period=10s --interval=5s --timeout=3s CMD ["./health-check"]
ENTRYPOINT ["./server"]
