FROM golang:1.18.3-alpine3.16 AS builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 go build -o app ./cmd/client

FROM scratch

COPY --from=builder /build/app /

ENTRYPOINT ["/app"]
