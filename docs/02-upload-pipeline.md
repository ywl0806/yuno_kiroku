# 📘 Upload Pipeline Design

## 1. Design Goals

| 목표                 | 내용                                                           |
| -------------------- | -------------------------------------------------------------- |
| **데이터 무결성**    | 처리 실패 시 media_item 레코드는 반드시 존재, 상태로 추적 가능 |
| **재처리 가능 구조** | failed 상태의 미디어는 상태 초기화로 재처리 가능               |
| **비용 효율**        | S3 직접 업로드 (API 서버 대역폭 절약)                          |
| **역할 분리**        | ML 추론(Python)과 DB/스토리지(Go) 명확히 분리                  |
| **동영상 지원**      | 이미지/동영상 자동 분기, ffmpeg 기반 재인코딩                  |

---

## 2. State Machine

### 상태 정의

| 코드 | 이름       | 설명                                                      |
| ---- | ---------- | --------------------------------------------------------- |
| `01` | pending    | Presigned URL 발급 완료, S3 업로드 대기 또는 완료         |
| `02` | processing | Resize Worker 수신, 처리 진행 중                          |
| `03` | completed  | 처리 완료 (이미지: 리사이즈 완료 / 동영상: 재인코딩 완료) |
| `04` | failed     | 처리 실패                                                 |

### 상태 전이 다이어그램

```mermaid
stateDiagram-v2
    direction LR

    [*] --> pending : Presigned URL 발급\n+ DB 레코드 생성

    pending --> processing : SQS resize 메시지 수신\n(resize-worker)

    processing --> completed : 리사이즈 완료 (이미지)\n재인코딩 완료 (동영상)
    processing --> failed : 처리 실패

    completed --> SQS_face_recognition : 이미지만
    SQS_face_recognition --> [*] : ai-batch → face-recognition-worker

    completed --> [*] : 동영상
    failed --> [*]
```

### 전이 조건

| 전이                 | 조건                                           | 담당 컴포넌트                |
| -------------------- | ---------------------------------------------- | ---------------------------- |
| 생성 → `01`          | Presigned URL 발급 요청                        | API Server                   |
| `01` → `02`          | SQS resize 메시지 수신                         | resize-worker                |
| `02` → `03` (이미지) | 리사이즈 완료 — **이 시점부터 사진 조회 가능** | resize-worker                |
| `02` → `03` (동영상) | 재인코딩 완료                                  | video-processing-worker      |
| `*` → `04`           | 처리 실패                                      | resize-worker / video-worker |

---

## 3. 이미지 업로드 시퀀스

```mermaid
sequenceDiagram
    participant C as 클라이언트
    participant API as API Server
    participant S3 as S3
    participant SQS_R as SQS (resize)
    participant RW as resize-worker
    participant SQS as SQS (face-recognition)
    participant AI as ai-batch (Python)
    participant FRW as face-recognition-worker (Go)
    participant DB as PostgreSQL

    C->>API: POST /presigned-url
    API->>DB: INSERT media_items(status=01) + media_files
    API->>S3: Presigned PUT URL 생성
    API-->>C: {media_item_id, presigned_url}

    C->>S3: PUT (직접 업로드)
    S3-->>C: 200 OK
    S3->>SQS_R: PutObject 이벤트 알림 (prefix=original/)

    SQS_R-->>RW: ReceiveMessage (Long Polling)
    RW->>DB: status=02
    RW->>S3: 원본 다운로드
    RW->>RW: EXIF 파싱
    RW->>S3: view / thumbnail 업로드
    RW->>DB: media_files INSERT, status=03 (조회 가능)
    RW->>SQS_R: DeleteMessage
    RW->>SQS: SendMessage {media_item_id, family_id, view_storage_key}

    SQS-->>AI: ReceiveMessage
    AI->>S3: view 이미지 다운로드
    AI->>AI: InsightFace 추론 → faces [{bbox, embedding}]
    AI->>FRW: subprocess stdin JSON

    FRW->>DB: advisory lock 획득
    FRW->>DB: pgvector 코사인 유사도 검색
    FRW->>DB: identity upsert + face_detections INSERT
    FRW->>S3: 얼굴 크롭 업로드 (512×512 WebP)
    FRW->>DB: identity_face_imgs INSERT

    C->>API: GET /upload-batch/status (폴링)
    API-->>C: {is_completed: true}
```

---

## 4. 동영상 업로드 시퀀스

```mermaid
sequenceDiagram
    participant C as 클라이언트
    participant S3 as S3
    participant SQS_R as SQS (resize)
    participant RW as resize-worker
    participant SQS as SQS (video-processing)
    participant VW as video-processing-worker
    participant DB as PostgreSQL

    C->>S3: PUT (직접 업로드)
    S3-->>C: 200 OK
    S3->>SQS_R: PutObject 이벤트 알림 (prefix=original/)

    SQS_R-->>RW: ReceiveMessage (Long Polling)
    RW->>RW: 동영상 확장자 감지 (mp4/mov/avi/mkv/webm/m4v)
    RW->>DB: status=02
    RW->>SQS_R: DeleteMessage
    RW->>SQS: SendMessage {media_item_id, family_id, original_key, file_name, mime_type}

    SQS-->>VW: ReceiveMessage
    VW->>S3: 원본 동영상 다운로드
    VW->>VW: ffmpeg 리사이즈 (최대 1280px, H.264)
    VW->>VW: ffmpeg 썸네일 추출 (1초, WebP)
    VW->>S3: thumbnail / video 업로드
    VW->>DB: media_files INSERT, status=03
```

---

## 5. DB 선 생성 전략

S3 업로드 전에 DB 레코드를 먼저 생성하여 고아 파일(orphan file)을 방지한다.

```
1. DB: media_items(status=01) + media_files(role=original) 생성
2. S3: Presigned PUT URL 발급
3. 클라이언트에 반환
4. 클라이언트: S3 직접 업로드
```

클라이언트가 업로드를 포기하면 `status=01` 레코드만 남는다.
정기 cleanup job으로 오래된 pending 레코드를 삭제한다.

---

## 6. SQS 큐 설계

### 큐 목록

| 큐 이름            | 발행자        | 소비자                       | 메시지 형식                                                              |
| ------------------ | ------------- | ---------------------------- | ------------------------------------------------------------------------ |
| `resize`           | S3 이벤트     | resize-worker (Go)           | S3 PutObject 이벤트 (key, bucket)                                        |
| `face-recognition` | resize-worker | ai-batch (Python)            | `{media_item_id, family_id, view_storage_key}`                           |
| `video-processing` | resize-worker | video-processing-worker (Go) | `{media_item_id, family_id, original_storage_key, file_name, mime_type}` |

### face-recognition 큐 메시지

```json
{
  "media_item_id": "uuid",
  "family_id": "uuid",
  "view_storage_key": "view/familyId/mediaItemId.webp"
}
```

### video-processing 큐 메시지

```json
{
  "media_item_id": "uuid",
  "family_id": "uuid",
  "original_storage_key": "original/familyId/mediaItemId.mp4",
  "file_name": "mediaItemId.mp4",
  "mime_type": "video/mp4"
}
```

---

## 7. AI 처리 — Python + Go 하이브리드

### 역할 분담

```
[ ai-batch (Python) ]                       [ face-recognition-worker (Go CLI) ]
  ─────────────────────────                   ────────────────────────────────────
  1. SQS 메시지 수신                           1. advisory lock 획득 (family 단위)
  2. view 이미지 S3 다운로드 (임시 파일)        2. pgvector 코사인 유사도 검색
  3. InsightFace 얼굴 감지 (bbox)              3. identity 매칭 / 신규 생성 (임계값 0.5)
  4. 512차원 ArcFace 임베딩 계산               4. face_detections INSERT (위치 + 임베딩)
  5. subprocess로 Go CLI 호출                  5. view 이미지에서 얼굴 크롭 (패딩 30%)
     stdin JSON:                               6. 512×512 리사이즈 → S3 identities/ 저장
       media_item_id, family_id,               7. identity_face_imgs INSERT
       view_image_path, faces
```

### subprocess 인터페이스

```python
# Python → Go CLI 호출
payload = json.dumps({
    "media_item_id": "uuid",
    "family_id": "uuid",
    "view_image_path": "/tmp/tmpXXXX.webp",
    "faces": [
        {
            "face_location": {"top": 10, "right": 100, "bottom": 110, "left": 0},
            "embedding": [0.12, -0.34, ...]  # 512차원
        }
    ]
})
proc = subprocess.run([FACE_RECOGNITION_WORKER], input=payload.encode(), capture_output=True)
```

---

## 8. Failure Handling

### resize-worker 실패

```
SQS resize → resize-worker
    ↓ 실패
media_items.status = '04' (failed)
DeleteMessage 하지 않음 → visibility timeout 후 자동 재시도
일정 횟수 실패 시 DLQ로 이동
```

### ai-batch 실패 (개별 job)

```python
# 처리 성공 시에만 SQS 메시지 삭제
success = process_job(s3_client, job)
if success:
    _delete_message(sqs, receipt_handle)
# 실패 시: visibility timeout 후 자동 재시도
```

### video-processing-worker 실패 (개별 job)

```go
// 처리 성공 시에만 DeleteMessage
if err := svc.ProcessVideo(ctx, params); err != nil {
    // DeleteMessage 하지 않음 → visibility timeout 후 재시도
    continue
}
deleteMessage(ctx, sqsClient, queueURL, msg.ReceiptHandle)
```

### ECS Task 중단 (Fargate Spot 선점)

```python
def _handle_signal(sig, frame):
    global _running
    logger.info("종료 신호 수신 - 현재 job 완료 후 종료")
    _running = False  # 루프를 gracefully 종료

signal.signal(signal.SIGTERM, _handle_signal)
signal.signal(signal.SIGINT, _handle_signal)
```

---

## 9. 컴포넌트 상세 스펙

### resize-worker (Go + govips)

```
런타임:  Go 1.24 + govips (libvips)
실행:    SQS long-polling (WaitTimeSeconds=20, prefix=original/)
처리:    이미지 → view(2048px) + thumbnail(512px) → SQS face-recognition
         동영상 → SQS video-processing
DI:      SQSFaceRecognitionDispatcher, SQSVideoJobDispatcher
```

### ai-batch (Python + InsightFace)

```
런타임:  Python 3.11 + InsightFace(buffalo_l) + boto3
메모리:  ~4-8 GB (ArcFace 모델 ~2GB + 배치 이미지)
실행:    SQS long-polling (WaitTimeSeconds=20, MaxMessages=10)
역할:    ML 추론만 수행 (얼굴 감지 + 임베딩)
후처리:  face-recognition-worker (Go CLI) subprocess 호출
```

### face-recognition-worker (Go CLI)

```
런타임:  Go 1.24 CLI 바이너리 (ai-batch Docker 이미지에 포함)
입력:    stdin JSON {media_item_id, family_id, view_image_path, faces}
처리:    advisory lock → identity 매칭 → face_detections INSERT → 얼굴 크롭 업로드
출력:    exit code 0 (성공) / 非0 (실패)
```

### video-processing-worker (Go + ffmpeg)

```
런타임:  Go 1.24 + ffmpeg-go (ffprobe, ffmpeg)
실행:    SQS long-polling (WaitTimeSeconds=20, VisibilityTimeout=900)
처리:    원본 다운로드 → ffmpeg 리사이즈(max 1280px, H.264) → 썸네일 추출(WebP) → S3 업로드
```

---

## 10. S3 버킷 구조

```
yuno-media-bucket/
├── original/     ← 클라이언트 Presigned URL 업로드 대상
│   └── {familyId}/{mediaItemId}.{jpg|heic|png|mp4|mov|...}
├── view/         ← resize-worker 출력 (2048px WebP)
│   └── {familyId}/{mediaItemId}.webp
├── thumbnail/    ← resize-worker 출력 (512px WebP) + 동영상 썸네일
│   └── {familyId}/{mediaItemId}.webp
├── video/        ← video-processing-worker 출력 (H.264 MP4)
│   └── {familyId}/{mediaItemId}.mp4
└── identities/   ← face-recognition-worker 출력 (얼굴 크롭 512×512)
    └── {identityId}/face_{mediaItemId}.webp
```

| prefix                | 접근 방법                              | 이유                 |
| --------------------- | -------------------------------------- | -------------------- |
| `original/`           | Presigned URL만 쓰기, 퍼블릭 읽기 불가 | 원본 보호            |
| `view/`, `thumbnail/` | CloudFront OAC로만 읽기                | 인증된 사용자만 접근 |
| `video/`              | CloudFront OAC로만 읽기                | 인증된 사용자만 접근 |
| `identities/`         | CloudFront OAC로만 읽기                | 얼굴 크롭 보호       |

---

## 11. 클라이언트 폴링 설계

### 폴링 엔드포인트

```
GET /media-item/upload-batch/status?upload_batch_id={id}

Response:
{
  "statuses": [
    { "media_item_id": "uuid", "upload_status": "03" },
    { "media_item_id": "uuid", "upload_status": "02" }
  ],
  "is_completed": false  // 모든 항목이 03/04/05일 때 true
}
```

### 단계별 예상 대기 시간

#### 이미지 파이프라인

| 단계                                    | 예상 시간       | 비고                                      |
| --------------------------------------- | --------------- | ----------------------------------------- |
| S3 업로드                               | 클라이언트 결정 | 파일 크기 / 네트워크                      |
| SQS resize → resize-worker 수신         | 수 초           | Long Polling (WaitTimeSeconds=20)         |
| 이미지 리사이즈 처리                    | 1-10초          | 원본 파일 크기                            |
| SQS → ai-batch 수신                     | 수 초           | Long Polling (WaitTimeSeconds=20)         |
| ↳ _CloudWatch SQS 적체 감지_            | _1-3분_         | _평가 기간 1분 × 연속 1-3회_              |
| ↳ _ECS ai-batch Task 스케일아웃 트리거_ | _수 초_         | _Application Auto Scaling 반응_           |
| ↳ _ECS Fargate 컨테이너 기동_           | _30-60초_       | _이미지 Pull 포함. 모델 로드(~10초) 추가_ |
| 얼굴 감지 + Go CLI 처리                 | 1-10초/장       | 얼굴 수, 해상도                           |
| **이미지 총 예상 (태스크 기동 전)**     | **~4-6분**      | **CloudWatch 감지 + 컨테이너 기동 포함**  |
| **이미지 총 예상 (태스크 가동 중)**     | **~10-30초**    | **얼굴 없는 사진은 리사이즈 직후 완료**   |

#### 동영상 파이프라인

| 단계                                              | 예상 시간       | 비고                                     |
| ------------------------------------------------- | --------------- | ---------------------------------------- |
| S3 업로드                                         | 클라이언트 결정 | 파일 크기 / 네트워크                     |
| SQS resize → resize-worker 수신                   | 수 초           | Long Polling (WaitTimeSeconds=20)        |
| SQS video-processing → Worker 수신                | 수 초           | Long Polling (WaitTimeSeconds=20)        |
| ↳ _CloudWatch SQS 적체 감지_                      | _1-3분_         | _평가 기간 1분 × 연속 1-3회_             |
| ↳ _ECS video-processing-worker 스케일아웃 트리거_ | _수 초_         | _Application Auto Scaling 반응_          |
| ↳ _ECS Fargate 컨테이너 기동_                     | _30-60초_       | _이미지 Pull 포함. ffmpeg 바이너리 포함_ |
| 동영상 재인코딩 (ffmpeg)                          | 30초~수분       | 영상 길이·해상도                         |
| **동영상 총 예상 (태스크 기동 전)**               | **~5-10분**     | **CloudWatch 감지 + 컨테이너 기동 포함** |
| **동영상 총 예상 (태스크 가동 중)**               | **~1-5분**      | **영상 길이에 비례**                     |

> _이탤릭_ 항목은 태스크가 0대일 때 처음 스케일아웃 되는 경우에만 발생하는 오버헤드입니다.  
> 태스크가 이미 가동 중이면 스케일아웃 없이 Long Polling 수신 즉시 처리됩니다.
