FROM python:3.11-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    libglib2.0-0 \
    && rm -rf /var/lib/apt/lists/* 

WORKDIR /app

# 먼저 requirements만 복사하여 캐시 활용
COPY requirements.txt .
RUN pip install --upgrade pip
RUN pip install -r requirements.txt

# 모델 디렉토리 복사 (로컬에 미리 다운로드된 모델)
COPY models/ /root/.insightface/models/

# 나머지 소스 코드 복사
COPY src/ ./src/

CMD ["python", "-m", "src.batch.local"]