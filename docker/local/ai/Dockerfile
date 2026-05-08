FROM golang:1.24-bookworm AS go-builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    libvips-dev libheif-dev pkg-config gcc g++ \
    && rm -rf /var/lib/apt/lists/*

ENV CGO_ENABLED=1
WORKDIR /src
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ .
RUN go build -o /face-recognition-worker ./cmd/face-recognition-worker

FROM python:3.11-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    libglib2.0-0 \
    libvips \
    libheif1 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY ai/requirements.txt .
RUN pip install --upgrade pip
RUN pip install -r requirements.txt

COPY ai/models/buffalo_l/ /root/.insightface/models/buffalo_l/
COPY ai/src/ ./src/
COPY --from=go-builder /face-recognition-worker ./face-recognition-worker

CMD ["python", "-m", "src.batch.sqs"]
