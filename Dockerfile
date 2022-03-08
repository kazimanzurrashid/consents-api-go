FROM golang:1.17.5-alpine3.15 AS builder
WORKDIR /usr/app
COPY . .
ENV CGO_ENABLED=0 \
    GOOS=linux
RUN go get -d -t ./... && go test ./...
RUN go build ./cmd/server && go build ./cmd/health-check

FROM scratch
ARG POSTGRES_HOST
ARG POSTGRES_USER
ARG POSTGRES_PASSWORD
ARG POSTGRES_DB
ARG PORT

ENV POSTGRES_HOST=$POSTGRES_HOST \
    POSTGRES_USER=$POSTGRES_USER \
    POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
    POSTGRES_DB=$POSTGRES_DB \
    PORT=$PORT

EXPOSE $PORT

WORKDIR /usr/app
COPY --from=builder /usr/app/server /usr/app/health-check /usr/app/schema.sql ./
HEALTHCHECK --start-period=10s --interval=15s --timeout=1s CMD ["./health-check"]
ENTRYPOINT ["./server"]
