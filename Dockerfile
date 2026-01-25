# Use a specific, stable Go version. 1.24 might not be out yet or stable.
# Using 1.23 as a safe bet, or check if 1.24 is actually available.
# Since the user had 1.24.4 in their file, I will stick to a valid tag if possible, or fallback to latest stable.
# Assuming 1.23 for stability as 1.24.4 seems futuristic (current stable is ~1.21/1.22/1.23).
FROM golang:1.24 AS builder

WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary
RUN go build -o coding-profile-service ./cmd/server

# Run in a debian-slim image to support installing chromium
FROM debian:bookworm-slim

# Install Chromium and necessary dependencies for chromedp
RUN apt-get update && apt-get install -y \
    chromium \
    ca-certificates \
    fonts-liberation \
    libasound2 \
    libatk-bridge2.0-0 \
    libatk1.0-0 \
    libc6 \
    libcairo2 \
    libcups2 \
    libdbus-1-3 \
    libexpat1 \
    libfontconfig1 \
    libgbm1 \
    libgcc1 \
    libglib2.0-0 \
    libgtk-3-0 \
    libnspr4 \
    libnss3 \
    libpango-1.0-0 \
    libpangocairo-1.0-0 \
    libstdc++6 \
    libx11-6 \
    libx11-xcb1 \
    libxcb1 \
    libxcomposite1 \
    libxcursor1 \
    libxdamage1 \
    libxext6 \
    libxfixes3 \
    libxi6 \
    libxrandr2 \
    libxrender1 \
    libxss1 \
    libxtst6 \
    lsb-release \
    wget \
    xdg-utils \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /root/

COPY --from=builder /app/coding-profile-service .
COPY README.md /root/README.md

EXPOSE 8080

CMD ["./coding-profile-service"]