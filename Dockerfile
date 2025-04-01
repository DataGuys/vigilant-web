  # Stage 1: Build the Go application
  FROM golang:1.18 AS builder
  WORKDIR /app
  COPY . .
  RUN go mod download
  RUN go build -o vigilant-onion ./cmd/main.go

  # Stage 2: Create a minimal image with the compiled binary
  FROM alpine:latest
  RUN apk --no-cache add ca-certificates
  WORKDIR /root/
  COPY --from=builder /app/vigilant-onion .
  CMD ["./vigilant-onion"]
