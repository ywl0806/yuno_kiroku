# AWS 인프라 구성도

> Phase 1 기준 (개인 운용 — VPC 없음, Supabase DB, ECS On-demand)

---

## 1. 전체 시스템 구성

```mermaid
graph TB
    subgraph Client["클라이언트"]
        WEB[Web Browser]
        MOB[Mobile App]
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
        L_RESIZE["Resize Lambda\nGo + libvips\n─────────────\nPOST /resize\nPOST /face-recognition/complete\nPOST /face-recognition/fail\n(Function URL 공개)"]
    end

    subgraph AI["ECS Fargate Spot"]
        ECS_AI["AI Task\nPython + InsightFace\n─────────────\n얼굴 감지 + 임베딩만 수행\n결과는 Function URL로 콜백"]
    end

    subgraph DB["Database"]
        SUPA[(Supabase\nPostgreSQL 17\n+ pgvector)]
    end

    subgraph Registry["ECR"]
        ECR_API[api-lambda image]
        ECR_RESIZE[resize-lambda image]
        ECR_AI[ai-task image]
    end

    %% 클라이언트 → 진입점
    WEB -->|HTTPS| CF_FRONT
    WEB -->|HTTPS| APIGW
    MOB -->|HTTPS| APIGW

    %% Frontend 정적 파일
    CF_FRONT -->|OAC| S3_FRONT

    %% API 흐름
    APIGW --> L_API
    L_API -->|쿼리| SUPA
    L_API -->|Presigned PUT URL 발급| S3_MEDIA

    %% 미디어 파일 서빙
    WEB -->|이미지 요청| CF_MEDIA
    MOB -->|이미지 요청| CF_MEDIA
    CF_MEDIA -->|OAC| S3_MEDIA

    %% 클라이언트 → S3 직접 업로드
    WEB -->|Presigned PUT| S3_MEDIA
    MOB -->|Presigned PUT| S3_MEDIA

    %% S3 이벤트 → Resize Lambda
    S3_MEDIA -->|"PutObject 이벤트\nprefix: original/"| L_RESIZE

    %% Resize Lambda 처리
    L_RESIZE -->|"view/(2048px)\nthumbnail/(512px) 저장"| S3_MEDIA
    L_RESIZE -->|"media_files INSERT\nupload_status=03"| SUPA
    L_RESIZE -->|"face_recognition_jobs INSERT\nECS ListTasks → RunTask"| ECS_AI

    %% ECS: 얼굴 감지만 수행
    ECS_AI -->|"pending jobs fetch"| SUPA
    ECS_AI -->|"view 이미지 다운로드"| S3_MEDIA
    ECS_AI -->|"POST /face-recognition/complete\n{ job_id, faces:[{location, embedding}] }\nFunction URL 콜백"| L_RESIZE

    %% Resize Lambda: 콜백 처리
    L_RESIZE -->|"identity 매칭 (pgvector)\nface_detections INSERT\nidentity_face_imgs INSERT\njob 완료"| SUPA
    L_RESIZE -->|"얼굴 크롭(512×512)\nidentities/ 저장"| S3_MEDIA

    %% ECR → Lambda/ECS 이미지 참조
    ECR_API -.->|image| L_API
    ECR_RESIZE -.->|image| L_RESIZE
    ECR_AI -.->|image| ECS_AI
```

---

## 2. 업로드 파이프라인 상세

```mermaid
sequenceDiagram
    actor C as 클라이언트
    participant APIGW as API Gateway
    participant API as API Lambda
    participant S3 as S3
    participant DB as Supabase DB
    participant RL as Resize Lambda
    participant ECS as ECS AI Task

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

    Note over S3,RL: S3 PutObject 이벤트 (prefix: original/)
    S3-)RL: invoke (POST /resize)
    RL->>DB: UPDATE media_items status=02 (processing)
    RL->>S3: 원본 다운로드
    RL->>RL: EXIF 파싱, 리사이즈
    RL->>S3: view/{familyId}/{mediaItemId}.webp 업로드 (2048px)
    RL->>S3: thumbnail/{familyId}/{mediaItemId}.webp 업로드 (512px)
    RL->>DB: INSERT media_files (view, thumbnail)
    RL->>DB: UPDATE media_items status=03 (completed) ← 사진 조회 가능
    RL->>DB: INSERT face_recognition_jobs (status=01)
    RL->>ECS: ListTasks → 미실행 시 RunTask (Fargate Spot)

    Note over ECS: 얼굴 감지만 수행 (임베딩 계산)
    ECS->>DB: fetch pending jobs (FOR UPDATE SKIP LOCKED)
    ECS->>S3: view 이미지 다운로드
    ECS->>ECS: InsightFace 얼굴 감지\n512차원 ArcFace 임베딩 계산
    ECS->>RL: POST /face-recognition/complete (Function URL)\n{ job_id, faces: [{ face_location, embedding }] }

    Note over RL: identity 매칭 + 크롭은 Resize Lambda가 처리
    RL->>DB: advisory lock 획득 (family 단위)
    RL->>DB: FindMostSimilarFace (pgvector 코사인 유사도)
    alt 유사 identity 없음
        RL->>DB: INSERT identities (새 인물 생성)
    end
    RL->>DB: INSERT face_detections
    RL->>S3: view 이미지 다운로드 (크롭용)
    RL->>RL: 얼굴 bbox + 패딩 크롭 → 512×512 리사이즈
    RL->>S3: identities/{identityId}/face_{mediaItemId}.webp 저장
    RL->>DB: INSERT identity_face_imgs
    RL->>DB: UPDATE face_recognition_jobs status=03 (completed)

    C->>APIGW: GET /api/media-item/upload-batch/status (폴링)
    APIGW->>API: invoke
    API->>DB: SELECT upload_status
    API-->>C: { is_completed: true }
```

---

## 3. 네트워크 토폴로지 (Phase 1 — VPC 없음)

```mermaid
graph TB
    subgraph Internet["인터넷"]
        USER[사용자]
        ECS_OUT[ECS → Function URL 콜백]
    end

    subgraph AWS_Region["AWS ap-northeast-1 (도쿄)"]
        subgraph Global["글로벌 / 엣지"]
            ACM_US[ACM 인증서\nus-east-1]
            CF[CloudFront\n엣지 POP]
        end

        subgraph NoVPC["VPC 없음 — 퍼블릭 인터넷 접근"]
            APIGW[API Gateway\nHTTP API]

            subgraph Lambda_Group["Lambda Functions"]
                L_API[API Lambda\n512MB / arm64\n타임아웃: 30s]
                L_RESIZE["Resize Lambda\n2048MB / x86_64\n타임아웃: 300s\n─────────────\nFunction URL (공개)\nPOST /resize\nPOST /face-recognition/complete\nPOST /face-recognition/fail"]
            end

            subgraph ECS_Group["ECS Fargate Spot"]
                ECS["AI Task\n2vCPU / 8GB\nassignPublicIp: ENABLED\n─────────────\n얼굴 감지 + 임베딩만 수행\n결과 → Function URL 콜백"]
            end

            subgraph Managed["관리형 서비스"]
                S3[S3\nyuno-media-bucket\nyuno-frontend]
                ECR[ECR\nDocker 이미지 저장소]
            end
        end
    end

    subgraph External["외부 서비스"]
        SUPA[Supabase\nPostgreSQL + pgvector\n공용 인터넷 TLS]
    end

    USER -->|HTTPS| CF
    CF -->|OAC| S3
    CF -->|Route| APIGW
    APIGW --> L_API
    L_API -->|인터넷 TLS| SUPA
    L_API -->|AWS SDK| S3
    L_RESIZE -->|인터넷 TLS| SUPA
    L_RESIZE -->|AWS SDK| S3
    L_RESIZE -->|ECS SDK\nRunTask| ECS
    ECS -->|인터넷 TLS| SUPA
    ECS -->|AWS SDK| S3
    ECS -->|"Function URL\nPOST /face-recognition/complete"| L_RESIZE
    ECR -.->|image pull| L_API
    ECR -.->|image pull| L_RESIZE
    ECR -.->|image pull| ECS
```

---

## 4. 역할 분담 요약

```mermaid
graph LR
    subgraph RL["Resize Lambda"]
        RL1[원본 다운로드]
        RL2[view / thumbnail 리사이즈]
        RL3[S3 업로드]
        RL4[media_files DB 저장]
        RL5[upload_status = 03\n사진 조회 가능]
        RL6[ECS RunTask 트리거]
        RL7["← Function URL 수신\n/face-recognition/complete"]
        RL8[identity 매칭\npgvector 유사도]
        RL9[face_detections INSERT]
        RL10[얼굴 크롭 + S3 저장]
        RL11[identity_face_imgs INSERT]
    end

    subgraph ECS["ECS AI Task (Python)"]
        ECS1[pending jobs fetch]
        ECS2[view 이미지 다운로드]
        ECS3[InsightFace 얼굴 감지]
        ECS4[512차원 ArcFace 임베딩]
        ECS5["Function URL 콜백\nPOST /face-recognition/complete\n{ job_id, faces }"]
    end

    RL6 -->|RunTask| ECS1
    ECS5 -->|HTTP 콜백| RL7
```

---

## 5. S3 버킷 구조

```mermaid
graph LR
    subgraph S3_MEDIA["S3: yuno-media-bucket"]
        direction TB
        ORIG["original/{familyId}/{mediaItemId}.jpg\n← Presigned PUT 업로드 대상\n← S3 이벤트로 Resize Lambda 트리거\n퍼블릭 읽기 불가"]
        VIEW["view/{familyId}/{mediaItemId}.webp\n← Resize Lambda 출력 (2048px)\n← ECS가 얼굴 감지용으로 다운로드\nCloudFront OAC만 읽기"]
        THUMB["thumbnail/{familyId}/{mediaItemId}.webp\n← Resize Lambda 출력 (512px)\nCloudFront OAC만 읽기"]
        IDENT["identities/{identityId}/face_{mediaItemId}.webp\n← Resize Lambda 콜백 처리 후 저장\n(얼굴 크롭 512×512)\nCloudFront OAC만 읽기"]
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
        R_RESIZE[resize-lambda-role]
        R_ECS[ecs-task-role]
        R_ECS_EXEC[ecs-task-execution-role]
    end

    subgraph Permissions["권한 (주요)"]
        P1["S3: GetObject, PutObject\nyuno-media-bucket"]
        P2["S3: GetObject, PutObject\nyuno-media-bucket"]
        P3["ECS: RunTask, ListTasks"]
        P4["S3: GetObject, PutObject\nyuno-media-bucket"]
        P5["ECR: GetAuthorizationToken\nBatchGetImage"]
        P6["Logs: CreateLogGroup\nPutLogEvents"]
    end

    R_API --> P1
    R_RESIZE --> P2
    R_RESIZE --> P3
    R_ECS --> P4
    R_ECS_EXEC --> P5
    R_ECS_EXEC --> P6
```

---

## 7. Phase 1 vs Phase 2 비교

```mermaid
graph LR
    subgraph P1["Phase 1 — 현재 (개인 운용)"]
        direction TB
        P1_DB[Supabase\n외부 PostgreSQL]
        P1_ECS[ECS Fargate Spot\nOn-demand Task]
        P1_NET[VPC 없음\nNAT Gateway 없음]
        P1_COST[비용: 사용량 기반\n유휴 시 과금 없음]
    end

    subgraph P2["Phase 2 — 미래 (실 서비스)"]
        direction TB
        P2_DB[RDS PostgreSQL 17\nVPC 내부]
        P2_ECS[ECS Fargate Service\n상시 가동 + Auto Scaling]
        P2_NET[VPC + NAT Gateway\nPrivate Subnet]
        P2_COST[비용: 고정 + 사용량\n안정성 우선]
    end

    P1 -->|"트리거: 사용자 증가\n트리거: 안정성 요구"| P2
```
