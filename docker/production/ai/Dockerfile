FROM python:3.11-slim AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    libglib2.0-0 \
    && rm -rf /var/lib/apt/lists/*

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# insightface buffalo_s 모델 사전 다운로드 (컨테이너 시작 시 딜레이 방지)
RUN python -c "\
from insightface.app import FaceAnalysis; \
app = FaceAnalysis(name='buffalo_s', providers=['CPUExecutionProvider']); \
app.prepare(ctx_id=0, det_size=(640, 640))"

FROM python:3.11-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    libgl1 \
    libglib2.0-0 \
    libgomp1 \
    curl \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages
COPY --from=builder /usr/local/bin /usr/local/bin
COPY --from=builder /root/.insightface /root/.insightface

WORKDIR /app
COPY src/ ./src/
RUN mkdir -p tmp/files

CMD ["python", "-m", "src.batch.sqs"]
