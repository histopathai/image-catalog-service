# Stage 1: Builder
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git bash

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o image-catalog-service ./cmd/main.go

# Stage 2: Final Image
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/image-catalog-service .

CMD ["./image-catalog-service"]
