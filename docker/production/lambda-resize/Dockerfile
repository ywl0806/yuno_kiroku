# ---------- Build stage ----------
FROM golang:1.24-bookworm AS builder

RUN apt-get update && apt-get install -y \
    libvips-dev libheif-dev \
    pkg-config gcc g++ \
    && rm -rf /var/lib/apt/lists/*

ENV CGO_ENABLED=1
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -tags lambda.norpc -o /main ./cmd/lambda-resize/main.go

# ---------- Runtime stage ----------
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    libvips libheif1 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /main /main

ENTRYPOINT ["/main"]