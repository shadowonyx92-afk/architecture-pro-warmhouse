# Stage 1: build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o device-service main.go

# Stage 2: run
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/device-service .
EXPOSE 8082
CMD ["./device-service"]

