# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o order_service ./cmd/server/order/main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/order_service .
COPY configs/order/.env ./configs/order/.env
EXPOSE 8082
CMD ["./order_service"]
