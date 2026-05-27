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

## 2. 업로드 파이프라인

> 이미지/동영상 업로드 시퀀스 상세: [02-upload-pipeline.md](02-upload-pipeline.md)
>
> 컴포넌트 역할 요약: [01-system-architecture-overview.md §2](01-system-architecture-overview.md)
>
> S3 버킷 구조: [02-upload-pipeline.md §10](02-upload-pipeline.md)

---

## 3. IAM 권한 구조

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
