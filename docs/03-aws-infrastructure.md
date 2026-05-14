# AWS 인프라 구성도

> Phase 1 기준 (개인 운용 — VPC 없음, Supabase DB, ECS On-demand)

---

## 1. 전체 시스템 구성

```mermaid
graph TB
    subgraph Client["클라이언트"]
        WEB[Web Browser]
    end

    subgraph CDN["CDN / 진입점"]
        CF_FRONT[CloudFront\nFrontend]
        CF_MEDIA[CloudFront\nMedia OAC]
        APIGW[API Gateway\nHTTP API]
    end

    subgraph Storage["S3"]
        S3_FRONT[S3\nyuno-frontend]
        S3_MEDIA[S3\nyuno-media-bucket]
    end

    subgraph Compute["Lambda"]
        L_API[API Lambda\nGo + Echo]
    end

    subgraph Workers["ECS Fargate Spot"]
        ECS_RESIZE["Resize Worker\nGo + govips\n─────────────\nHTTP 서버 (Echo)\nPOST /resize\n(MinIO Webhook 수신)\n→ 이미지: SQS face-recognition 발행\n→ 동영상: SQS video-processing 발행"]
        ECS_AI["AI Batch\nPython + InsightFace\n─────────────\nSQS face-recognition 폴링\n얼굴 감지 + 임베딩\n→ subprocess: face-recognition-worker"]
        ECS_VIDEO["Video Worker\nGo + ffmpeg\n─────────────\nSQS video-processing 폴링\n동영상 리사이즈 (1280px H.264)\n썸네일 추출 (WebP)"]
    end

    subgraph Queues["SQS"]
        SQS_FACE[face-recognition\n큐]
        SQS_VIDEO[video-processing\n큐]
    end

    subgraph DB["Database"]
        SUPA[(Supabase\nPostgreSQL 17\n+ pgvector)]
    end

    subgraph Registry["ECR"]
        ECR_API[api-lambda image]
        ECR_RESIZE[resize-worker image]
        ECR_AI[ai-batch image\n+ face-recognition-worker 바이너리 포함]
        ECR_VIDEO[video-worker image]
    end

    %% 클라이언트 → 진입점
    WEB -->|HTTPS| CF_FRONT
    WEB -->|HTTPS| APIGW

    %% Frontend 정적 파일
    CF_FRONT -->|OAC| S3_FRONT

    %% API 흐름
    APIGW --> L_API
    L_API -->|쿼리| SUPA
    L_API -->|Presigned PUT URL 발급| S3_MEDIA

    %% 미디어 파일 서빙
    WEB -->|이미지/동영상 요청| CF_MEDIA
    CF_MEDIA -->|OAC| S3_MEDIA

    %% 클라이언트 → S3 직접 업로드
    WEB -->|Presigned PUT| S3_MEDIA

    %% S3 이벤트 → Resize Worker
    S3_MEDIA -->|"PutObject Webhook\nprefix: original/"| ECS_RESIZE

    %% Resize Worker 처리
    ECS_RESIZE -->|"view/(2048px)\nthumbnail/(512px) 저장"| S3_MEDIA
    ECS_RESIZE -->|"media_files INSERT\nupload_status=03"| SUPA
    ECS_RESIZE -->|"face-recognition job 발행"| SQS_FACE
    ECS_RESIZE -->|"video-processing job 발행"| SQS_VIDEO

    %% AI Batch: SQS 폴링 → 얼굴 감지 → Go CLI
    SQS_FACE -->|"SQS ReceiveMessage"| ECS_AI
    ECS_AI -->|"view 이미지 다운로드"| S3_MEDIA
    ECS_AI -->|"subprocess: face-recognition-worker\n{media_item_id, family_id,\n view_image_path, faces}"| ECS_AI
    ECS_AI -->|"identity 매칭 (pgvector)\nface_detections INSERT\nidentity_face_imgs INSERT"| SUPA
    ECS_AI -->|"얼굴 크롭(512×512)\nidentities/ 저장"| S3_MEDIA

    %% Video Worker: SQS 폴링 → 동영상 처리
    SQS_VIDEO -->|"SQS ReceiveMessage"| ECS_VIDEO
    ECS_VIDEO -->|"원본 동영상 다운로드"| S3_MEDIA
    ECS_VIDEO -->|"thumbnail/video 업로드\nmedia_files INSERT\nupload_status=03"| SUPA
    ECS_VIDEO -->|"video/(1280px H.264)\nthumbnail/(WebP) 저장"| S3_MEDIA

    %% ECR → 이미지 참조
    ECR_API -.->|image| L_API
    ECR_RESIZE -.->|image| ECS_RESIZE
    ECR_AI -.->|image| ECS_AI
    ECR_VIDEO -.->|image| ECS_VIDEO
```

---

## 2. 이미지 업로드 파이프라인 상세

```mermaid
sequenceDiagram
    actor C as 클라이언트
    participant APIGW as API Gateway
    participant API as API Lambda
    participant S3 as S3
    participant DB as Supabase DB
    participant RW as Resize Worker
    participant SQS_F as SQS\nface-recognition
    participant AI as AI Batch\n(Python)
    participant GoCLI as face-recognition\nworker (Go CLI)

    C->>APIGW: POST /api/media-item/upload-batch
    APIGW->>API: invoke
    API->>DB: INSERT upload_batch
    API-->>C: { batch_id }

    C->>APIGW: POST /api/media-item/presigned-url
    APIGW->>API: invoke
    API->>DB: INSERT media_items(status=01) + media_files
    API->>S3: GeneratePresignedPutURL
    API-->>C: { media_item_id, presigned_url }

    C->>S3: PUT original/{familyId}/{mediaItemId}.jpg
    S3-->>C: 200 OK

    Note over S3,RW: MinIO/S3 PutObject Webhook (prefix: original/)
    S3-)RW: POST /resize
    RW->>DB: UPDATE media_items status=02 (processing)
    RW->>S3: 원본 다운로드
    RW->>RW: EXIF 파싱, 리사이즈
    RW->>S3: view/{familyId}/{mediaItemId}.webp 업로드 (2048px)
    RW->>S3: thumbnail/{familyId}/{mediaItemId}.webp 업로드 (512px)
    RW->>DB: INSERT media_files (view, thumbnail)
    RW->>DB: UPDATE media_items status=03 (completed) ← 사진 조회 가능
    RW->>SQS_F: SendMessage {media_item_id, family_id, view_storage_key}

    Note over AI: SQS Long Polling
    AI->>SQS_F: ReceiveMessage
    AI->>S3: view 이미지 다운로드 (임시 파일)
    AI->>AI: InsightFace 얼굴 감지\n512차원 ArcFace 임베딩 계산
    AI->>GoCLI: subprocess stdin:\n{media_item_id, family_id,\n view_image_path, faces}

    Note over GoCLI: identity 매칭 + 크롭
    GoCLI->>DB: advisory lock 획득 (family 단위)
    GoCLI->>DB: FindMostSimilarFace (pgvector 코사인 유사도)
    alt 유사 identity 없음
        GoCLI->>DB: INSERT identities (새 인물 생성)
    end
    GoCLI->>DB: INSERT face_detections
    GoCLI->>GoCLI: 얼굴 bbox + 패딩 크롭 → 512×512 리사이즈
    GoCLI->>S3: identities/{identityId}/face_{mediaItemId}.webp 저장
    GoCLI->>DB: INSERT identity_face_imgs
    AI->>SQS_F: DeleteMessage (처리 완료)

    C->>APIGW: GET /api/media-item/upload-batch/status (폴링)
    APIGW->>API: invoke
    API->>DB: SELECT upload_status
    API-->>C: { is_completed: true }
```

---

## 3. 동영상 업로드 파이프라인 상세

```mermaid
sequenceDiagram
    actor C as 클라이언트
    participant S3 as S3
    participant DB as Supabase DB
    participant RW as Resize Worker
    participant SQS_V as SQS\nvideo-processing
    participant VW as Video Worker\n(Go + ffmpeg)

    C->>S3: PUT original/{familyId}/{mediaItemId}.mp4
    S3-->>C: 200 OK

    Note over S3,RW: MinIO/S3 PutObject Webhook
    S3-)RW: POST /resize
    RW->>RW: 동영상 확장자 감지 (mp4/mov/avi/...)
    RW->>DB: UPDATE media_items status=02 (processing)
    RW->>SQS_V: SendMessage {media_item_id, family_id, original_key, ...}

    Note over VW: SQS Long Polling (VisibilityTimeout=900s)
    VW->>SQS_V: ReceiveMessage
    VW->>S3: 원본 동영상 스트리밍 다운로드
    VW->>VW: ffmpeg 리사이즈\n(max 1280px, H.264 CRF23, faststart)
    VW->>VW: ffmpeg 썸네일 추출\n(1초 지점, WebP)
    VW->>S3: thumbnail/{familyId}/{mediaItemId}.webp 저장
    VW->>S3: video/{familyId}/{mediaItemId}.mp4 저장
    VW->>DB: INSERT media_files (thumbnail + video)
    VW->>DB: UPDATE media_items status=03 (completed)
    VW->>SQS_V: DeleteMessage
```

---

## 4. 역할 분담 요약

```mermaid
graph LR
    subgraph RW["Resize Worker (Go)"]
        RW1[Webhook 수신\nPOST /resize]
        RW2[원본 다운로드]
        RW3[view / thumbnail 리사이즈]
        RW4[S3 업로드]
        RW5[media_files DB 저장]
        RW6[upload_status = 03\n사진 조회 가능]
        RW7[SQS face-recognition 발행]
        RW8[SQS video-processing 발행]
    end

    subgraph AI["AI Batch (Python)"]
        AI1[SQS 폴링]
        AI2[view 이미지 다운로드]
        AI3[InsightFace 얼굴 감지]
        AI4[512차원 ArcFace 임베딩]
        AI5[subprocess: Go CLI 호출]
    end

    subgraph GoCLI["face-recognition-worker (Go CLI)"]
        GO1[advisory lock 획득]
        GO2[pgvector identity 매칭]
        GO3[face_detections INSERT]
        GO4[얼굴 크롭 + S3 저장]
        GO5[identity_face_imgs INSERT]
    end

    subgraph VW["Video Worker (Go + ffmpeg)"]
        VW1[SQS 폴링]
        VW2[원본 동영상 다운로드]
        VW3[ffmpeg 리사이즈\nH.264 max 1280px]
        VW4[ffmpeg 썸네일 추출\nWebP]
        VW5[S3 업로드]
        VW6[upload_status = 03]
    end

    RW7 -->|SQS 메시지| AI1
    AI5 -->|subprocess stdin| GO1
    RW8 -->|SQS 메시지| VW1
```

---

## 5. S3 버킷 구조

```mermaid
graph LR
    subgraph S3_MEDIA["S3: yuno-media-bucket"]
        direction TB
        ORIG["original/{familyId}/{mediaItemId}.{ext}\n← Presigned PUT 업로드 대상\n← Webhook으로 Resize Worker 트리거\n퍼블릭 읽기 불가"]
        VIEW["view/{familyId}/{mediaItemId}.webp\n← Resize Worker 출력 (2048px)\nCloudFront OAC만 읽기"]
        THUMB["thumbnail/{familyId}/{mediaItemId}.webp\n← Resize Worker 출력 (이미지 512px)\n← Video Worker 출력 (동영상 썸네일)\nCloudFront OAC만 읽기"]
        VIDEO["video/{familyId}/{mediaItemId}.mp4\n← Video Worker 출력 (H.264 1280px)\nCloudFront OAC만 읽기"]
        IDENT["identities/{identityId}/face_{mediaItemId}.webp\n← face-recognition-worker 출력\n(얼굴 크롭 512×512)\nCloudFront OAC만 읽기"]
    end

    subgraph S3_FRONT["S3: yuno-frontend"]
        direction TB
        STATIC["index.html / assets/\n← React 빌드 결과물\nCloudFront OAC만 읽기"]
    end
```

---

## 6. IAM 권한 구조

```mermaid
graph LR
    subgraph Roles["IAM Roles"]
        R_API[api-lambda-role]
        R_RESIZE[resize-worker-role]
        R_AI[ai-batch-role]
        R_VIDEO[video-worker-role]
        R_ECS_EXEC[ecs-task-execution-role]
    end

    subgraph Permissions["권한 (주요)"]
        P_API["S3: GetObject, PutObject\nyuno-media-bucket"]
        P_RESIZE["S3: GetObject, PutObject\nyuno-media-bucket\nSQS: SendMessage\nface-recognition, video-processing"]
        P_AI["S3: GetObject, PutObject\nyuno-media-bucket\nSQS: ReceiveMessage, DeleteMessage\nface-recognition"]
        P_VIDEO["S3: GetObject, PutObject\nyuno-media-bucket\nSQS: ReceiveMessage, DeleteMessage\nvideo-processing"]
        P_EXEC["ECR: GetAuthorizationToken\nBatchGetImage\nLogs: PutLogEvents"]
    end

    R_API --> P_API
    R_RESIZE --> P_RESIZE
    R_AI --> P_AI
    R_VIDEO --> P_VIDEO
    R_ECS_EXEC --> P_EXEC
```

---
