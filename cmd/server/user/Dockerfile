# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o user_service ./cmd/server/user/main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/user_service .
COPY configs/user/.env ./configs/user/.env
EXPOSE 8080
CMD ["./user_service"]
