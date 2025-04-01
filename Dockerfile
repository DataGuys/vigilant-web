  # Use the official Golang image to create a build artifact
  FROM golang:1.18 AS builder
  WORKDIR /app
  COPY . .
  RUN go mod download
  RUN go build -o main .

  # Use a minimal image to run the binary
  FROM alpine:latest
  RUN apk --no-cache add ca-certificates
  WORKDIR /root/
  COPY --from=builder /app/main .
  CMD ["./main"]
