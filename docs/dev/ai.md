# AI 서비스 (얼굴인식)

**디렉토리:** `ai/`

---

## 개요

ArcFace 모델을 사용한 얼굴 감지·얼굴 임베딩 벡터 생성 서비스입니다. SQS 큐에서 얼굴인식 job을 폴링하여 배치 처리합니다. 업로드된 사진에서 얼굴을 감지하고 512차원의 임베딩 벡터를 생성합니다. 이 벡터를 pgvector로 검색하여 동일인물 사진을 자동으로 그룹핑합니다.

### 로컬 개발과 프로덕션의 차이

| 환경 | 방식 |
|------|------|
| 로컬 개발 | ElasticMQ (SQS 호환) 폴링 |
| 프로덕션 (Phase 1) | ECS Fargate Spot + Supabase DB 폴링 |
| 프로덕션 (Phase 2) | ECS Fargate Service + SQS 폴링 |

---

## 사용 기술

| 카테고리 | 라이브러리 | 버전 |
|----------|-----------|------|
| 얼굴인식 모델 | InsightFace | 0.7.3 |
| 이미지 처리 | OpenCV (headless) | 4.12.0 |
| 이미지 처리 | Pillow | 11.3.0 |
| 데이터 검증 | Pydantic | 2.12.4 |
| 수치 계산 | NumPy | 2.0.2 |
| 머신러닝 | scikit-learn | 1.6.1 |

---

## 디렉토리 구조

```
ai/
├── src/
│   ├── batch/
│   │   └── sqs.py               # SQS 폴링 배치 진입점
│   ├── service/
│   │   └── face.py              # 얼굴 감지 로직
│   └── utils/                   # 유틸리티 함수
├── models/                      # 사전 학습된 모델 파일
├── requirements.txt             # Python 의존성
└── docker/                      # Docker 설정
```

---

## 배치 처리 흐름

```
SQS 큐 (face-recognition)
    ↓ 폴링 (src/batch/sqs.py)
메시지 수신 (media_item_id 포함)
    ↓
S3/MinIO에서 view 이미지 다운로드
    ↓
InsightFace (buffalo_s) 얼굴 감지
    ↓
512차원 ArcFace 임베딩 계산
    ↓
pgvector 코사인 유사도로 identity 매칭/신규 생성 (임계값 0.5)
    ↓
face_detections DB 저장
    ↓
얼굴 크롭 (512×512) → S3/MinIO identities/ 저장
    ↓
face_recognition_jobs 완료 처리
```

---

## 얼굴인식 모델

### InsightFace (ArcFace)

- **모델:** `buffalo_s`
- **출력:** 512차원 정규화된 얼굴 임베딩 벡터
- **유사도 판정:** 코사인 유사도 (임계값 0.5)

### 처리 파이프라인

```
입력 이미지
    ↓
OpenCV / Pillow 전처리
    ↓
InsightFace (ArcFace) 얼굴 감지
    ↓
바운딩 박스 좌표 추출
    ↓
각 얼굴의 512차원 임베딩 생성
    ↓
정규화 후 pgvector에 저장
```

---

## 환경 변수

| 변수 | 설명 |
|------|------|
| `DATABASE_URL` | PostgreSQL 연결 URL |
| `STORAGE_TYPE` | `minio` 또는 `s3` |
| `S3_ENDPOINT` | MinIO 엔드포인트 (로컬: `http://minio:9000`) |
| `MEDIA_BUCKET_NAME` | 버킷 이름 (`my-bucket`) |
| `SQS_QUEUE_URL` | face-recognition 큐 URL |
| `SQS_ENDPOINT_URL` | SQS 엔드포인트 (로컬: `http://queue:9324`) |

---

## 개발 명령어

```bash
# Docker로 실행 (권장)
make ai-batch
# 또는
docker compose up -d ai-batch

# 직접 실행
cd ai
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python -m src.batch.sqs
```

---

## 주의 사항

- **GPU 지원:** InsightFace는 CPU/GPU 모두 지원하나, 로컬 개발은 CPU 모드로 동작합니다.
- **모델 파일:** `models/` 디렉토리에 사전 학습된 모델이 필요합니다. 초회 기동 시 자동 다운로드됩니다.
- **메모리:** ArcFace 모델이 약 2GB 메모리를 사용합니다. 컨테이너 메모리 제한에 주의하세요.
