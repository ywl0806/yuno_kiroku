# 📘 Upload Pipeline Design

## 1. Design Goals

| 목표                 | 내용                                                              |
| -------------------- | ----------------------------------------------------------------- |
| **데이터 무결성**    | 처리 실패 시 media_item 레코드는 반드시 존재, 상태로 추적 가능    |
| **Idempotency**      | 동일 요청 중복 처리 시 DB에 중복 레코드 생성되지 않음             |
| **재처리 가능 구조** | failed 상태의 job은 attempt_count 기반으로 자동 재시도            |
| **비용 효율**        | S3 직접 업로드 (API 서버 대역폭 절약), ECS는 job 존재 시에만 기동 |
| **Lambda 호환**      | 고루틴 수명 보장 불필요, 각 처리 단계가 독립 Lambda/Task로 분리   |
| **단계적 확장**      | Phase 1 (On-demand ECS) → Phase 2 (상시 ECS) 전환 시 최소 변경    |

---

## 2. State Machine

### 상태 정의

| 코드 | 이름       | 설명                                                  |
| ---- | ---------- | ----------------------------------------------------- |
| `01` | pending    | Presigned URL 발급 완료, S3 업로드 대기 또는 완료     |
| `02` | processing | Resize Lambda 처리 완료, ECS AI 처리 대기/진행 중     |
| `03` | completed  | ECS AI 처리 완료 (얼굴인식, identity 매칭, 크롭 저장) |
| `04` | failed     | 처리 실패 (재처리 대상)                               |
| `05` | duplicate  | 중복 사진으로 판정 (벡터 유사도 기반)                 |

### 상태 전이 다이어그램

```
              ┌──────────────────────────┐
              │  Presigned URL 발급      │
              │  + DB 레코드 생성        │
              └────────────┬─────────────┘
                           │
                           ▼
                      ┌─────────┐
                      │   01    │  pending
                      └────┬────┘
                           │  S3 PutObject 이벤트 → Resize Lambda
                           ▼
                      ┌─────────┐
                      │   02    │  processing
                      │         │  (리사이즈 완료 + SQS signal)
                      └────┬────┘
                           │  ECS AI 처리
               ┌───────────┼───────────┐
               ▼           ▼           ▼
          ┌─────────┐ ┌────────┐ ┌─────────┐
          │   03    │ │   04   │ │   05    │
          │completed│ │ failed │ │duplicate│
          └─────────┘ └───┬────┘ └─────────┘
                          │  attempt_count < 3
                          ▼
                     재시도 (01로 복구)
```

### 전이 조건

| 전이        | 조건                                        | 담당 컴포넌트               |
| ----------- | ------------------------------------------- | --------------------------- |
| 생성 → `01` | Presigned URL 발급 요청                     | API Lambda                  |
| `01` → `02` | Resize Lambda 시작 (S3 이벤트 수신)         | Resize Lambda               |
| `02` → `03` | 리사이즈 완료 — **이 시점부터 사진 조회 가능** | Resize Lambda             |
| `03` → `04` | 처리 실패                                   | Resize Lambda               |
| ─           | (face_recognition_jobs가 별도 상태 추적)    | Resize Lambda (콜백 처리)   |

---

## 3. Upload Sequence

### 전체 시퀀스 다이어그램

```
클라이언트     API Lambda      S3        Resize Lambda              ECS AI (Python)
    │              │            │               │                          │
    │ POST /batch  │            │               │                          │
    │─────────────▶│            │               │                          │
    │ { batch_id } │            │               │                          │
    │◀─────────────│            │               │                          │
    │              │            │               │                          │
    │ POST /presigned-url       │               │                          │
    │─────────────▶│            │               │                          │
    │              │ INSERT     │               │                          │
    │              │ media_items│               │                          │
    │              │ media_files│               │                          │
    │              │ PresignPut │               │                          │
    │              │───────────▶│               │                          │
    │ { media_item_id,          │               │                          │
    │   presigned_url }         │               │                          │
    │◀─────────────│            │               │                          │
    │              │            │               │                          │
    │ PUT (직접 업로드)         │               │                          │
    │───────────────────────────▶               │                          │
    │ 200 OK       │            │               │                          │
    │◀───────────────────────────               │                          │
    │              │            │               │                          │
    │              │            │ PutObject 이벤트 (prefix: original/)     │
    │              │            │──────────────▶│                          │
    │              │            │◀──────────────│ 원본 다운로드             │
    │              │            │◀──────────────│ view/thumbnail 업로드    │
    │              │            │               │ UPDATE status=02→03      │
    │              │            │               │ INSERT face_recognition_jobs
    │              │            │               │ ECS ListTasks            │
    │              │            │               │─────────────────────────▶│ RunTask
    │              │            │               │                          │ (미실행 시)
    │              │            │               │                          │
    │              │            │               │                  job fetch (DB 폴링)
    │              │            │◀─────────────────────────────────────────│ view 다운로드
    │              │            │               │                          │ 얼굴 감지 + 임베딩
    │              │            │               │◀─────────────────────────│
    │              │            │               │  POST /face-recognition/complete (Function URL)
    │              │            │               │  { job_id, faces: [{face_location, embedding}] }
    │              │            │               │                          │
    │              │            │               │ identity 매칭 (pgvector) │
    │              │            │               │ face_detections INSERT   │
    │              │            │◀──────────────│ 얼굴 크롭 업로드 (identities/)
    │              │            │               │ identity_face_imgs INSERT│
    │              │            │               │ job status=completed     │
    │              │            │               │                          │ job 없으면
    │              │            │               │                          │ Task 자동 종료
    │ GET /batch/status (폴링)  │               │                          │
    │─────────────▶│            │               │                          │
    │ { is_completed: true }    │               │                          │
    │◀─────────────│            │               │                          │
```

### DB 선 생성 전략

S3 업로드 전에 DB 레코드를 먼저 생성하여 고아 파일(orphan file)을 방지한다.

```
1. DB: media_items (status=01) + media_files (role=01) 생성
2. S3: Presigned PUT URL 발급
3. 클라이언트에 반환
4. 클라이언트: S3 직접 업로드
```

클라이언트가 업로드를 포기하면 `status=01` 레코드만 남는다.
24시간 후 cleanup job으로 삭제한다.

---

## 4. ECS 직접 트리거 (Phase 1)

### 설계 원칙

- **SQS / trigger-ecs Lambda 없음** — Resize Lambda가 ECS RunTask를 직접 호출
- **ECS는 얼굴 감지만** — 임베딩 계산 후 결과를 Resize Lambda Function URL로 콜백
- **후처리는 Resize Lambda** — identity 매칭, face_detections INSERT, 얼굴 크롭, DB 업데이트
- ECS가 이미 실행 중이면 RunTask 스킵 (ListTasks로 확인)
- ECS Task는 DB에서 pending job을 모두 처리하고 스스로 종료

```
Resize Lambda (ECSFaceRecognitionDispatcher)
    └─ DB: face_recognition_jobs INSERT
    └─ ECS: ListTasks → running_count > 0이면 스킵
    └─ ECS: RunTask (Fargate Spot) → 미실행 시에만

ECS AI Task (Python) — 얼굴 감지만 수행
    └─ while True:
    │      jobs = fetch_pending_jobs(limit=50)
    │      if not jobs: break  ← pending 없으면 Task 종료
    │      for job in jobs:
    │          image = s3.download(job.view_storage_key)
    │          faces = insightface.detect(image)  ← 감지 + 임베딩만
    │          POST {RESIZE_LAMBDA_FUNCTION_URL}/face-recognition/complete
    │               { job_id, faces: [{ face_location, embedding }] }
    └─ sys.exit(0)  ← 자연 종료 → ECS Task 완료 상태

Resize Lambda — Function URL 콜백 (/face-recognition/complete)
    └─ goroutine으로 비동기 처리:
    │      for face in faces:
    │          advisory lock (family 단위)
    │          pgvector 코사인 유사도 → identity 매칭/생성
    │          face_detections INSERT
    │          view 이미지 크롭(512×512) → S3 identities/ 저장
    │          identity_face_imgs INSERT
    │      face_recognition_jobs status=completed
    └─ 즉시 200 OK 반환 (ECS 대기 없이)
```

### Phase 2 전환 시 변경

Phase 2에서는 ECS Service가 상시 가동되므로 RunTask 직접 호출 방식이 변경된다.

```
Phase 1: Resize Lambda → ECS RunTask (On-demand) → Function URL 콜백
Phase 2: ECS Service 상시 가동 (SQS 폴링) → Function URL 콜백 (유지)
```

---

## 5. Idempotency Strategy

### upload_batch_id 전략

클라이언트가 업로드 세션 시작 시 발급받는 batch_id.
같은 배치 내 동일 파일명 중복 요청 방지:

```sql
UNIQUE (upload_batch_id, file_name)
```

### 얼굴 기반 중복 이미지 검사

pgvector 코사인 유사도로 거의 동일한 사진을 감지한다.

```python
# ECS AI Task: 얼굴 임베딩이 1개 이상인 사진에서만 중복 검사
def check_duplicate(conn, family_id, faces):
    for face in faces:
        row = find_most_similar_face(conn, family_id, face['embedding'])
        if row and row['distance'] < DUPLICATE_THRESHOLD:  # 0.05 이하
            return True
    return False
```

### SQS 중복 메시지 처리

SQS 메시지가 중복 수신되더라도 face_recognition_jobs 테이블의 UNIQUE 제약으로 중복 job 생성을 방지한다.
trigger-ecs Lambda는 ECS running 여부만 확인하므로 중복 수신 시 불필요한 ECS 중복 기동이 발생하지 않는다.

```sql
CONSTRAINT uq_frj_media_item_id UNIQUE (media_item_id)
```

---

## 6. Failure Handling

### 업로드 중단 (클라이언트 취소)

- S3 이벤트 미발생 → Resize Lambda 미실행
- DB에 `status=01` 레코드만 잔류
- 24시간 후 cleanup job이 삭제
- 클라이언트 재시도: 새 Presigned URL로 재업로드 가능

### Resize Lambda 실패

```
S3 이벤트 → Resize Lambda 실행
    ↓ 실패
Lambda 내장 재시도 (최대 2회)
    ↓ 3회 모두 실패
SQS DLQ (resize-failed-dlq)
    ↓
CloudWatch Alarm → SNS → 알림
```

```go
if err := processResize(ctx, key, mediaItemID); err != nil {
    db.UpdateMediaItemUploadStatus(ctx, mediaItemID, "04")
    return err  // Lambda 재시도 후 DLQ
}
```

### ECS AI Task 실패 (개별 job)

```python
def handle_failure(conn, job_id, media_item_id, attempt_count, error_msg):
    if attempt_count >= MAX_RETRIES:  # 3
        # 완전 실패
        UPDATE face_recognition_jobs SET status='04', last_error=...
        UPDATE media_items SET upload_status='04'
    else:
        # 다음 ECS 기동 시 재처리 (status를 01로 복구)
        UPDATE face_recognition_jobs SET status='01', last_error=...
```

### ECS Task 중단 (Fargate Spot 선점)

Fargate Spot은 2분 전 SIGTERM 신호 후 Task 종료.

```python
def sigterm_handler(signum, frame):
    """처리 중인 job을 pending으로 롤백"""
    conn.execute("""
        UPDATE face_recognition_jobs SET status = '01'
        WHERE status = '02'
    """)
    conn.commit()
    sys.exit(0)

signal.signal(signal.SIGTERM, sigterm_handler)
```

---

## 7. Reprocessing Strategy

### failed 상태 재처리

`attempt_count < 3`: ECS 다음 기동 시 자동 재처리 (status='01'로 자동 복구).
`attempt_count >= 3`: 수동 초기화 후 재처리.

```sql
UPDATE face_recognition_jobs
SET status = '01', attempt_count = 0, last_error = NULL
WHERE status = '04' AND media_item_id = ?;

UPDATE media_items SET upload_status = '01' WHERE id = ?;
```

### processing stuck 회수

trigger-ecs Lambda 실행 시, 또는 ECS Task 시작 직후 다음 쿼리를 실행한다.

```sql
-- 30분 이상 processing 상태인 job → pending으로 복구
UPDATE face_recognition_jobs
SET status = '01'
WHERE status = '02'
  AND scheduled_at < NOW() - INTERVAL '30 minutes';
```

---

## 8. Data Consistency Guarantees

### 고아 파일 방지

```
전략 1: DB 선 생성
  - S3 업로드 전 media_items + media_files 레코드 먼저 생성
  - S3 실패 시 DB 레코드 존재 → cleanup 가능

전략 2: S3 Object Tag
  - Presigned URL 생성 시 x-amz-tagging: media_item_id=1234 포함
  - Resize Lambda에서 GetObjectTagging으로 media_item_id 즉시 조회
  - media_files.storage_key 기반 DB 조회 불필요

전략 3: S3 Lifecycle 정책
  - original/: 90일 후 Glacier 이전 (비용 절감)
  - view/, thumbnail/: 유지 (서빙용)
```

### 순서 보장

S3 이벤트와 SQS 메시지는 순서가 보장되지 않는다.
Resize Lambda는 각 S3 Object마다 독립 실행되므로 순서 의존성이 없다.
ECS AI 배치는 `created_at ASC`로 처리하되, 클라이언트 폴링은 배치 단위로 완료 여부를 확인하므로 순서 무관하게 동작한다.

---

## 9. 컴포넌트 상세 스펙

### API Lambda (Go)

```
런타임:  Go 1.23 + Echo + algnhsa 어댑터
메모리:  512 MB
타임아웃: 30초
아키텍처: arm64 (Graviton2, 비용 20% 절감)
VPC:    없음 (Phase 1) / 있음 (Phase 2)
Provisioned Concurrency: 2 (콜드 스타트 방지)
배포:   ECR Docker 컨테이너
```

### Resize Lambda (Go + libvips)

```
런타임:  Go 1.23 + govips (CGO 필수)
메모리:  2048 MB
타임아웃: 300초
아키텍처: x86_64 (libvips CGO 호환성)
VPC:    없음 (Phase 1) / 있음 (Phase 2)
Reserved Concurrency: 10 (Supabase 커넥션 풀 보호, Phase 1)
트리거:  S3 PutObject, prefix=original/
DLQ:   SQS resize-failed-dlq (최종 실패 시)
배포:   ECR Docker 컨테이너 (libvips 포함)
```

```dockerfile
# Dockerfile.resize-lambda
FROM golang:1.23 AS builder
RUN apt-get update && apt-get install -y libvips-dev
WORKDIR /app
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -o bootstrap ./cmd/resize-lambda/

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y libvips42 && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/bootstrap /var/runtime/bootstrap
ENTRYPOINT ["/var/runtime/bootstrap"]
```

### ECS Fargate AI Task (Python)

```
런타임:  Python 3.11 + InsightFace (buffalo_s) + psycopg2 + boto3
CPU:    2 vCPU
메모리:  8 GB (ArcFace 모델 ~2GB + 배치 이미지)
실행방식 (Phase 1): Fargate Spot, On-demand Task (job 없으면 자동 종료)
실행방식 (Phase 2): Fargate Service, 상시 가동 + SQS Auto Scaling
배치크기: 50 jobs/실행

역할: 얼굴 감지 + 임베딩 계산만 수행
  → 결과를 Resize Lambda Function URL로 콜백
  → identity 매칭, DB 업데이트, 얼굴 크롭은 Resize Lambda가 처리

콜백 엔드포인트:
  POST {RESIZE_LAMBDA_FUNCTION_URL}/face-recognition/complete
    { job_id: int, faces: [{ face_location: {top,right,bottom,left}, embedding: [512]float }] }
  POST {RESIZE_LAMBDA_FUNCTION_URL}/face-recognition/fail
    { job_id: int, error: string }
```

### Supabase (Phase 1 DB)

```
서비스:  Supabase PostgreSQL 17 + pgvector
접속:   공용 인터넷 (TLS), DATABASE_URL 환경변수
장점:   VPC 불필요, 무료 tier 가능, pgvector 기본 지원
전환:   Phase 2에서 RDS로 교체 (pg_dump → RDS restore)
```

---

## 10. face_recognition_jobs 테이블 스키마

```sql
CREATE TABLE face_recognition_jobs (
    id               SERIAL PRIMARY KEY,
    media_item_id    INTEGER NOT NULL REFERENCES media_items(id) ON DELETE CASCADE,
    family_id        INTEGER NOT NULL REFERENCES families(id),
    view_storage_key VARCHAR(512) NOT NULL,      -- ECS가 다운로드할 S3 경로
    status           VARCHAR(2) NOT NULL DEFAULT '01',
    -- '01': pending / '02': processing / '03': completed / '04': failed
    attempt_count    INTEGER NOT NULL DEFAULT 0,
    last_error       TEXT,
    scheduled_at     TIMESTAMP,                 -- ECS가 fetch한 시각
    completed_at     TIMESTAMP,
    created_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_frj_media_item_id UNIQUE (media_item_id)
);

-- 배치 fetch 성능
CREATE INDEX idx_frj_status_created_at
    ON face_recognition_jobs (status, created_at) WHERE status = '01';

-- stuck 회수 쿼리 성능
CREATE INDEX idx_frj_status_scheduled_at
    ON face_recognition_jobs (status, scheduled_at) WHERE status = '02';

-- media_item_id 조회
CREATE INDEX idx_frj_media_item_id ON face_recognition_jobs (media_item_id);
```

---

## 11. S3 버킷 구조

```
yuno-media-bucket/
├── original/     ← 클라이언트 Presigned URL 업로드 대상
│   └── {familyId}/{albumId}/{YYYY-MM-DD}/{uuid}.{jpg|heic|png}
├── view/         ← Resize Lambda 출력 (2048px WebP)
│   └── {familyId}/{albumId}/{YYYY-MM-DD}/{uuid}.webp
├── thumbnail/    ← Resize Lambda 출력 (512px WebP)
│   └── {familyId}/{albumId}/{YYYY-MM-DD}/{uuid}.webp
└── identities/   ← ECS AI 출력 (얼굴 크롭 512×512 WebP)
    └── {identityId}/face_{mediaItemId}.webp
```

| prefix                | 접근 방법                              | 이유                 |
| --------------------- | -------------------------------------- | -------------------- |
| `original/`           | Presigned URL만 쓰기, 퍼블릭 읽기 불가 | 원본 보호            |
| `view/`, `thumbnail/` | CloudFront OAC로만 읽기                | 인증된 사용자만 접근 |
| `identities/`         | CloudFront OAC로만 읽기                | 얼굴 크롭 보호       |

---

## 12. 클라이언트 폴링 설계

### 폴링 엔드포인트

```
GET /media-item/upload-batch/status?upload_batch_id={id}

Response:
{
  "statuses": [
    { "media_item_id": 1234, "upload_status": "03" },
    { "media_item_id": 1235, "upload_status": "02" },
    { "media_item_id": 1236, "upload_status": "04" }
  ],
  "is_completed": false  // 모든 항목이 03/04/05일 때 true
}
```

### 권장 폴링 구현 (Exponential Backoff)

```typescript
async function pollBatchStatus(uploadBatchId: number) {
  let interval = 3_000 // 시작: 3초 (Resize Lambda 처리 대기)
  const maxInterval = 10_000 // 최대: 10초 (ECS 기동 대기)

  while (true) {
    const result = await fetchBatchStatus(uploadBatchId)
    if (result.isCompleted) break

    interval = Math.min(interval * 1.2, maxInterval)
    await sleep(interval)
  }
}
```

### 단계별 예상 대기 시간

| 단계                        | 예상 시간       | 비고                                 |
| --------------------------- | --------------- | ------------------------------------ |
| S3 업로드                   | 클라이언트 결정 | 파일 크기 / 네트워크                 |
| S3 이벤트 → Resize Lambda   | 1-5초           |                                      |
| Resize Lambda 처리          | 5-30초          | 원본 파일 크기                       |
| SQS → trigger-ecs Lambda    | 수 초           |                                      |
| ECS Task 기동 (콜드 스타트) | 30-60초         | 컨테이너 이미지 pull                 |
| ECS AI 처리                 | 1-10초/장       | 얼굴 수, 해상도                      |
| **총 예상 시간**            | **1-3분**       | Spot 사용 시 Task 기동 시간이 지배적 |

---
