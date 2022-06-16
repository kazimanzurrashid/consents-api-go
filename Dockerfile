FROM golang:1.18.3-alpine3.16 AS builder
WORKDIR /usr/app
COPY . .
ENV CGO_ENABLED=0 \
    GOOS=linux
RUN go get -d -t ./... && go test ./...
RUN go build ./cmd/server && go build ./cmd/health-check

FROM scratch
ARG POSTGRES_HOST
ARG POSTGRES_PORT
ARG POSTGRES_USER
ARG POSTGRES_PASSWORD
ARG POSTGRES_DB
ARG PORT

ENV POSTGRES_HOST=$POSTGRES_HOST \
    POSTGRES_PORT=$POSTGRES_PORT \
    POSTGRES_USER=$POSTGRES_USER \
    POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
    POSTGRES_DB=$POSTGRES_DB \
    PORT=$PORT

EXPOSE $PORT

WORKDIR /usr/app
COPY --from=builder /usr/app/server /usr/app/health-check /usr/app/schema.sql ./
HEALTHCHECK --start-period=30s --interval=15s --timeout=1s CMD ["./health-check"]
ENTRYPOINT ["./server"]
