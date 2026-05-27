FROM golang:1.24-bookworm

RUN apt-get update && apt-get install -y --no-install-recommends \
    ffmpeg \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

RUN go install github.com/air-verse/air@v1.61.5

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air", "-c", ".air.video-worker.toml"]
