# 로컬 개발 환경 설정

---

## 사전 준비

| 도구 | 버전 | 용도 |
|------|------|------|
| Docker + Docker Compose | 최신 | 로컬 인프라 전체 |
| Go | 1.24.0+ | 백엔드 개발 |
| Node.js | 18+ | 프론트엔드·모바일 |
| Python | 3.10+ | AI 서비스 |

---

## 1. Docker 서비스 시작

프로젝트 루트에서 실행합니다.

```bash
# 전체 서비스 시작
make up

# 또는 직접
docker compose up -d
docker compose logs -f app resize-worker ai-batch
```

### 구동 서비스

| 서비스 | 포트 | 설명 |
|--------|------|------|
| `app` | 1323 | Go 백엔드 서버 |
| `resize-worker` | 1325 | 이미지 리사이즈 워커 (MinIO Webhook 수신) |
| `ai-batch` | - | AI 얼굴인식 배치 (SQS 폴링) |
| `video-processing-worker` | - | 동영상 처리 워커 (SQS 폴링) |
| `postgres` | 5433 | PostgreSQL 17 (pgvector) |
| `minio` | 9001 / 9090 | S3 호환 오브젝트 스토리지 |
| `queue` | 9324 / 9325 | ElasticMQ (SQS 호환) |

> MinIO 콘솔: `http://localhost:9090` (root / password)

> `app`과 `ai-batch`만 따로 시작하려면:
> ```bash
> docker compose up -d postgres minio queue
> ```

---

## 2. DB 초기화

```bash
# 마이그레이션 실행
make migrate

# 시드 데이터 삽입
make seed

# 한 번에 (destroy + migrate + seed)
make refresh

# DB 초기화만
make destroy
```

---

## 3. Web 프론트엔드 설정

```bash
cd front

# 의존성 설치
npm install

# 개발 서버 시작
npm run dev
```

Web 앱: `http://localhost:5155`

Vite proxy 설정:
- `/api` → `http://127.0.0.1:1323`
- `/uploads` → `http://127.0.0.1:1323`

---

## 4. 모바일 앱 설정

```bash
cd mobile

# 의존성 설치
npm install

# iOS: Pod 설치
cd ios && pod install && cd ..

# Metro 개발 서버 시작
npm start

# iOS 실행 (별도 터미널)
npm run ios

# Android 실행
npm run android
```

---

## 5. AI 배치 서비스 (Docker 미사용 시)

AI 서비스는 기본적으로 docker-compose의 `ai-batch`로 실행됩니다.
직접 실행하려면:

```bash
cd ai

# 가상환경 생성
python -m venv venv
source venv/bin/activate

# 의존성 설치
pip install -r requirements.txt

# 배치 실행 (SQS 폴링 모드)
python -m src.batch.sqs
```

---

## 환경 변수

### server/.env (또는 docker-compose 환경변수 참고)

```env
# 데이터베이스
DATABASE_URL=postgres://postgres:postgres@localhost:5432/yuno?sslmode=disable

# JWT
JWT_SECRET=your-secret-key

# OAuth - LINE
LINE_CLIENT_ID=your-line-client-id
LINE_CLIENT_SECRET=your-line-client-secret
LINE_REDIRECT_URL=http://localhost:1323/auth/callback/line

# OAuth - Kakao
KAKAO_CLIENT_ID=your-kakao-client-id
KAKAO_REDIRECT_URL=http://localhost:1323/auth/callback/kakao

# 스토리지 (MinIO)
STORAGE_TYPE=minio
S3_ENDPOINT=http://minio:9000
MINIO_ROOT_USER=root
MINIO_ROOT_PASSWORD=password
MEDIA_BUCKET_NAME=my-bucket

# SQS (ElasticMQ)
SQS_ENDPOINT_URL=http://queue:9324
VIDEO_SQS_QUEUE_URL=http://queue:9324/000000000000/video-processing
SQS_QUEUE_URL=http://queue:9324/000000000000/face-recognition
```

---

## 자주 쓰는 Make 명령어

```bash
# 서비스 시작 + 로그 확인
make up

# 서비스 중지
make stop

# 재시작 (app, resize-worker, ai-batch)
make reload

# DB 관련
make migrate    # 마이그레이션 실행
make seed       # 시드 데이터 삽입
make destroy    # DB 초기화
make refresh    # destroy + migrate + seed

# 코드 생성
make sqlc       # sqlc로 DB 쿼리 코드 생성
make swag       # Swagger 문서 생성

# 포트 정리
make kill       # 1323 포트 프로세스 종료
```

---

## Swagger UI

서버 시작 후 API 문서 확인:

```
http://localhost:1323/swagger/index.html
```

---

## 트러블슈팅

### PostgreSQL 연결 오류

```bash
docker compose ps
docker compose logs postgres
```

### pgvector 확장이 없는 경우

```sql
CREATE EXTENSION IF NOT EXISTS vector;
```

### AI 서비스가 시작되지 않는 경우

모델 파일이 `ai/models/`에 있는지 확인하세요. 없으면 초회 시작 시 자동 다운로드됩니다.

### MinIO Webhook 이벤트가 오지 않는 경우

`create-buckets` 컨테이너 로그를 확인하세요:

```bash
docker compose logs create-buckets
```

### 포트 충돌

`docker-compose.yml`에서 포트 번호를 변경하고, `.env` 파일도 맞게 수정하세요.
