FROM python:3.11-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    libglib2.0-0 \
    && rm -rf /var/lib/apt/lists/* 

WORKDIR /app

COPY . .
RUN pip install --upgrade pip
RUN pip install -r requirements.txt

# insightface 모델 사전 다운로드
RUN python -c "from insightface.app import FaceAnalysis; app = FaceAnalysis(name='buffalo_s', providers=['CPUExecutionProvider']); app.prepare(ctx_id=0, det_size=(640, 640))"

EXPOSE 8000

CMD ["uvicorn", "src.main:app", "--host", "0.0.0.0", "--reload", "--port", "8000"]