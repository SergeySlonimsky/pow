FROM golang:1.18.3-alpine3.16 AS builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 go build -o app ./cmd/server

FROM scratch

ARG SERVER_PORT=8083

COPY --from=builder /build/app /

EXPOSE ${SERVER_PORT}

ENTRYPOINT ["/app"]
