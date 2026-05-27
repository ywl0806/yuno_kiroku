# 백엔드 서버

**디렉토리:** `server/`

---

## 사용 기술

| 카테고리 | 라이브러리 | 버전 |
|----------|-----------|------|
| 언어 | Go | 1.24.0 |
| Web 프레임워크 | Echo | 4.11.4 |
| DB | PostgreSQL | 17 |
| 벡터 검색 | pgvector | - |
| SQL 코드 생성 | sqlc | - |
| DB 마이그레이션 | golang-migrate | 4.19.0 |
| 이미지 처리 | govips (libvips) | v2 |
| 유효성 검사 | go-playground/validator | v10 |
| JWT | golang-jwt | v5.3.0 |
| AWS 클라이언트 | aws-sdk-go-v2 | - |
| UUID | google/uuid | - |
| 설정 관리 | viper | - |
| EXIF | goexif | - |
| API 문서 | swaggo (Swagger) | - |

---

## 디렉토리 구조

```
server/
├── cmd/
│   ├── api/                    # Go API 서버 진입점
│   ├── resize/                 # 리사이즈 워커 진입점
│   ├── video-processing-worker/ # 동영상 처리 워커
│   ├── face-recognition-worker/ # 얼굴인식 워커
│   ├── lambda/                 # AWS Lambda 어댑터
│   └── db/                     # DB 유틸 (migrate, seed, destroy)
├── internal/
│   ├── api/
│   │   ├── app.go              # Echo 앱 초기화
│   │   ├── handlers/           # HTTP 요청/응답 처리
│   │   ├── routers/            # API 라우트 정의
│   │   └── middlewares/        # HTTP 미들웨어
│   ├── services/               # 비즈니스 로직
│   ├── store/                  # 데이터 접근 레이어
│   ├── stores/                 # 스토리지 추상화
│   ├── db/                     # sqlc 생성 코드 (직접 수정 금지)
│   ├── oauth/                  # OAuth 프로바이더 (LINE, Kakao)
│   ├── worker/                 # 워커 공통 로직
│   ├── apperr/                 # 커스텀 에러 타입
│   ├── consts/                 # 상수
│   ├── enums/                  # 열거형
│   ├── i18n/                   # 국제화 (에러 메시지 등)
│   ├── providers/              # 의존성 주입
│   └── validator/              # 커스텀 유효성 검사
├── db/
│   ├── schema.sql              # DB 스키마 정의
│   ├── queries/                # sqlc용 SQL 쿼리
│   └── migrations/             # 마이그레이션 파일
└── go.mod
```

---

## 아키텍처

Layered Architecture를 채택합니다.

```
HTTP 요청
    ↓
Middleware (CORS / 인증 / 로그)
    ↓
Handler (요청 파싱·응답 포맷)
    ↓
Service (비즈니스 로직)
    ↓
Store (데이터 접근)
    ↓
DB (PostgreSQL / sqlc 생성 코드)
```

### 워커 구조

이미지/동영상 처리는 별도 워커 프로세스로 분리되어 있습니다.

```
MinIO Webhook (PUT original/)
    ↓
resize-worker (Go + govips)
    리사이즈 (view 2048px, thumbnail 512px)
    ↓ SQS: face-recognition 큐
ai-batch (Python)
    얼굴 감지 + 임베딩 계산
    ↓ DB: face_detections, identities

동영상 업로드
    ↓ SQS: video-processing 큐
video-processing-worker (Go)
    동영상 트랜스코딩 처리
```

---

## API 엔드포인트

### 인증

| 메서드 | 경로 | 설명 |
|--------|------|------|
| GET | `/auth/callback/line` | LINE OAuth 콜백 |
| GET | `/auth/callback/kakao` | Kakao OAuth 콜백 |
| POST | `/auth/refresh` | 토큰 갱신 |
| POST | `/auth/logout` | 로그아웃 |

### 유저

| 메서드 | 경로 | 설명 |
|--------|------|------|
| GET | `/users/me` | 프로필 조회 |
| PUT | `/users/me` | 프로필 수정 |

### 미디어 아이템

| 메서드 | 경로 | 설명 |
|--------|------|------|
| POST | `/media-item/upload-batch` | 업로드 배치 생성 |
| POST | `/media-item/presigned-url` | Presigned URL 발급 |
| GET | `/media-item` | 날짜 범위로 조회 |
| GET | `/media-item/upload-batch/status` | 배치 업로드 상태 조회 |

### 앨범

| 메서드 | 경로 | 설명 |
|--------|------|------|
| GET | `/albums` | 앨범 목록 |
| POST | `/albums` | 앨범 생성 |
| GET | `/albums/:id` | 앨범 상세 |
| PUT | `/albums/:id` | 앨범 수정 |

### 인물(Identity)

| 메서드 | 경로 | 설명 |
|--------|------|------|
| GET | `/identities` | 인물 목록 |
| POST | `/identities` | 인물 생성 |
| PUT | `/identities/:id` | 인물 수정 |

### 그룹·멤버·아이

| 메서드 | 경로 | 설명 |
|--------|------|------|
| GET/POST/PUT | `/groups` | 그룹 관리 |
| GET/POST/PUT | `/kids` | 아이 관리 |

### 초대

| 메서드 | 경로 | 설명 |
|--------|------|------|
| POST | `/invites` | 초대 토큰 생성 |
| GET | `/invites/:token` | 초대 토큰 검증 |

### 좋아요·태그·설정

| 메서드 | 경로 | 설명 |
|--------|------|------|
| POST/DELETE | `/likes` | 좋아요 토글 |
| GET/POST/DELETE | `/tags` | 태그 관리 |
| GET/PUT | `/settings` | 앱 설정 |

---

## 데이터베이스 스키마

### 주요 테이블

| 테이블 | 설명 |
|--------|------|
| `families` | 가족 조직 |
| `groups` | 가족 내 그룹 |
| `users` | 유저 계정 (OAuth 연동) |
| `albums` | 사진 앨범 (그룹 권한 포함) |
| `media_items` | 사진·동영상 아이템 |
| `media_files` | 파일 실체 (original/thumbnail/view) |
| `identities` | 얼굴인식으로 식별된 인물 |
| `face_detections` | 얼굴 감지 데이터 (512차원 벡터) |
| `face_recognition_jobs` | 얼굴인식 처리 작업 큐 |
| `invite_tokens` | 초대 링크 토큰 |
| `refresh_tokens` | JWT 리프레시 토큰 |

### pgvector 얼굴 유사도 검색

```sql
-- 512차원 임베딩 저장
embedding vector(512)

-- IVFFlat 인덱스로 코사인 유사도 검색 고속화
CREATE INDEX ON face_detections
  USING ivfflat (embedding vector_cosine_ops);
```

### media_files 파일 역할

| 역할 | 설명 |
|------|------|
| `original` | 원본 고해상도 파일 |
| `thumbnail` | 목록 표시용 썸네일 (512px) |
| `view` | 뷰어용 최적화 파일 (2048px) |

---

## 인증·인가

- **JWT:** Access Token(HTTP-only Cookie) + Refresh Token(DB 저장)
- **OAuth 프로바이더:** LINE / Kakao
- **미들웨어:** 인증이 필요한 라우트에 `AuthMiddleware` 적용

---

## 파일 스토리지

- **로컬 개발:** MinIO (S3 호환, 포트 9001)
- **프로덕션:** AWS S3
- **버킷:** `my-bucket` (단일 버킷, prefix로 구분)
  - `original/{familyId}/` — 원본
  - `view/{familyId}/` — 뷰어용
  - `thumbnail/{familyId}/` — 썸네일
  - `identities/{identityId}/` — 얼굴 크롭
- **환경변수:** `STORAGE_TYPE=minio|s3|local`

---

## 개발 명령어

```bash
# DB 작업 (프로젝트 루트의 Makefile 사용)
make migrate    # 마이그레이션 실행
make seed       # 시드 데이터 삽입
make destroy    # DB 초기화
make refresh    # destroy + migrate + seed

# 코드 생성
make sqlc       # sqlc로 DB 쿼리 코드 생성
make swag       # Swagger 문서 생성

# 서비스 재시작
make reload     # app, resize-worker, ai-batch 재시작
```

---

## 환경변수

`server/.env.example`을 참고하여 `.env` 파일을 생성하세요.

Swagger UI: `http://localhost:1323/swagger/index.html`
