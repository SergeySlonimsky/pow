FROM golang:1.18.3-alpine3.16 AS builder

WORKDIR /build

COPY . .

RUN go build -o app ./cmd/client

FROM alpine:3.16

COPY --from=builder /build/app /

ENTRYPOINT ["/app"]
