FROM python:3.11-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    libgl1 \
    libglib2.0-0 \
    libgomp1 \
    curl \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# insightface buffalo_s 모델 사전 다운로드 (컨테이너 시작 시 딜레이 방지)
RUN python -c "\
from insightface.app import FaceAnalysis; \
app = FaceAnalysis(name='buffalo_s', providers=['CPUExecutionProvider']); \
app.prepare(ctx_id=0, det_size=(640, 640))"

COPY src/ ./src/

RUN mkdir -p tmp/files

CMD ["python", "-m", "src.batch.face_recognition_batch"]