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

FROM python:3.11-slim AS python-builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    libglib2.0-0 \
    && rm -rf /var/lib/apt/lists/*

COPY ai/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# RUN python -c "\
# from insightface.app import FaceAnalysis; \
# app = FaceAnalysis(name='buffalo_l', providers=['CPUExecutionProvider']); \
# app.prepare(ctx_id=0, det_size=(640, 640))"
COPY ai/models/buffalo_l/ /root/.insightface/models/buffalo_l/

FROM python:3.11-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    libgl1 \
    libglib2.0-0 \
    libgomp1 \
    libvips \
    libheif1 \
    curl \
    && rm -rf /var/lib/apt/lists/*

COPY --from=python-builder /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages
COPY --from=python-builder /usr/local/bin /usr/local/bin
COPY --from=python-builder /root/.insightface /root/.insightface

WORKDIR /app
COPY ai/src/ ./src/
COPY --from=go-builder /face-recognition-worker ./face-recognition-worker

RUN mkdir -p tmp/files

CMD ["python", "-m", "src.batch.sqs"]
