FROM golang:1.18.3-alpine3.16 AS builder

WORKDIR /build

COPY . .

RUN go build -o app ./cmd/server

FROM alpine:3.16

ARG SERVER_PORT=8083

COPY --from=builder /build/app /

EXPOSE ${SERVER_PORT}

ENTRYPOINT ["/app"]
