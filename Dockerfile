FROM golang:1.20 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o vigilant-web ./cmd/main.go

# final stage
FROM alpine:3.17
WORKDIR /app

# If you need CA certs for crawling, install them
RUN apk --no-cache add ca-certificates && update-ca-certificates

COPY --from=builder /app/vigilant-web /usr/local/bin/vigilant-web
COPY web/ /app/web/

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/vigilant-web"]
