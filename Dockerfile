# Use the latest Go image (>= 1.24)
FROM golang:1.24.4 AS builder

WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary
RUN go build -o coding-profile-service ./cmd/server

# Run in a lightweight image
FROM gcr.io/distroless/base-debian12

WORKDIR /root/

COPY --from=builder /app/coding-profile-service .
# ✅ Add this line:
COPY README.md /root/README.md

EXPOSE 8080

CMD ["./coding-profile-service"]