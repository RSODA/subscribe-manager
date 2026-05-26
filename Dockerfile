FROM golang:1.26.3-alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o subscribe-manager ./cmd/main

FROM alpine:3.23

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/subscribe-manager .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/.env ./.env

CMD ["./subscribe-manager"]
