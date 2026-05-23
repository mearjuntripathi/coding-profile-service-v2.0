# Stage 1: Build
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o coding-profile-service ./cmd/server

# Stage 2: Run - use minimal alpine image, no Chromium needed
FROM alpine:latest

# Just need CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/coding-profile-service .
COPY README.md /root/README.md

EXPOSE 8080

CMD ["./coding-profile-service"]