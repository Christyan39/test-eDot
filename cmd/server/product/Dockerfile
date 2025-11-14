# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o product_service ./cmd/server/product/main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/product_service .
COPY configs/product/.env ./configs/product/.env
EXPOSE 8081
CMD ["./product_service"]
