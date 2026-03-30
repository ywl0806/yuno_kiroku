# 📘 System Architecture Overview

## 1. Introduction

### 프로젝트 목적

YUNO는 가족 단위의 사진을 공유하고, 얼굴 인식 기반으로 자동 분류하는 사진 관리 플랫폼이다.
단순한 앨범 공유를 넘어, 가족 구성원의 얼굴을 자동으로 인식하고 인물별 검색 및 필터링을 제공한다.

### 대상 사용자

- **가족 단위**: 부모/자녀가 같은 Family에 속해 사진을 공유
- **멀티 디바이스**: Web, iOS, Android에서 동일한 경험 제공
- **비기술 사용자**: 업로드 후 자동 처리, 별도 조작 없이 얼굴 분류 완료

### 설계 목표

| 목표                       | 내용                                                   |
| -------------------------- | ------------------------------------------------------ |
| **비용 최소화**            | 유휴 시간 비용 제로 수준의 Serverless 기반 설계        |
| **단계적 확장**            | 개인 운용 → 실 서비스 전환 시 컴포넌트 교체만으로 대응 |
| **클라우드 종속성 최소화** | 인터페이스 추상화로 AWS → 홈서버 이전 가능한 구조      |
| **데이터 무결성**          | 이미지 처리 실패 시 재처리 가능, 고아 파일 방지        |

---

## 2. 아키텍처 단계

### Phase 1 — 개인 운용 (현재)

**목표**: 나 혼자 사용, 비용 최소화. 트래픽이 없을 때 과금 없음.

```
[ 클라이언트 ]
      │ HTTPS
      ▼
[ API Gateway + Go API Lambda ]
      │  DB: Supabase PostgreSQL (외부, VPC 불필요)
      │  S3: Presigned URL 발급
      ▼
[ S3: yuno-media-bucket ]
  original/ prefix → PutObject 이벤트
      │
      ▼
[ Resize Lambda (Go + libvips) ]
  - 리사이즈: view(2048px), thumbnail(512px)
  - DB: face_recognition_jobs INSERT
  - SQS: SendMessage (wake-up signal)
      │
      ▼
[ SQS: face-recognition-queue ]
      │  Lambda Trigger (SQS Event Source)
      ▼
[ trigger-ecs Lambda ]
  - ECS 실행 중 여부 확인
  - 실행 중 → 종료 (ECS가 알아서 DB 폴링 중)
  - 미실행 → ECS RunTask
      │
      ▼
[ ECS Fargate Spot Task: Python AI ]
  - Supabase에서 pending job 배치 fetch
  - 얼굴인식 → identity 매칭 → DB 저장
  - 더 이상 job 없으면 Task 종료 (비용 절감)
      │
      ▼
[ Supabase PostgreSQL + pgvector ]
```

**비용 포인트:**

- Lambda: 요청 없으면 과금 없음
- ECS: 처리할 job이 있을 때만 Task 기동, 완료 후 자동 종료
- EventBridge Scheduler 없음 (SQS 이벤트 기반으로 즉각 트리거)
- VPC 없음 → NAT Gateway 비용 없음 (Supabase는 공용 인터넷으로 접근)

---

### Phase 2 — 실 서비스 (미래)

**목표**: 다수 사용자 대응, 안정성 우선. ECS 상시 가동 + Auto Scaling.

```
[ 클라이언트 ]
      │ HTTPS
      ▼
[ API Gateway + Go API Lambda ]
      │  RDS PostgreSQL (VPC 내부)
      │  S3: Presigned URL 발급
      ▼
[ S3: yuno-media-bucket ]
      │ S3 이벤트 → Resize Lambda
      ▼
[ Resize Lambda (Go + libvips) ]
  - 리사이즈 완료 후 SQS SendMessage
      │
      ▼
[ SQS: face-recognition-queue ]
      │  ECS Service가 상시 폴링
      ▼
[ ECS Fargate Service: Python AI (상시 가동) ]
  - SQS ReceiveMessage 폴링
  - queue depth 기반 Auto Scaling
  - RDS PostgreSQL + pgvector
```

**Phase 1 → Phase 2 전환 시 변경 사항:**

| 컴포넌트           | Phase 1                   | Phase 2                   |
| ------------------ | ------------------------- | ------------------------- |
| Database           | Supabase (외부)           | RDS PostgreSQL (VPC 내부) |
| ECS 실행 방식      | On-demand (job 있을 때만) | 상시 가동 Service         |
| SQS 소비자         | ECS Task (DB 폴링)        | ECS Service (SQS 폴링)    |
| trigger-ecs Lambda | 필요                      | 불필요 (ECS가 상시 대기)  |
| VPC / NAT          | 불필요                    | 필요                      |
| 비용 구조          | 사용량 기반               | 고정 + 사용량 혼합        |

---

## 3. 주요 컴포넌트

| 컴포넌트          | Phase 1                          | Phase 2                      |
| ----------------- | -------------------------------- | ---------------------------- |
| **API Server**    | Lambda + API Gateway             | 동일                         |
| **Resize Worker** | Lambda (Go + libvips)            | 동일                         |
| **AI Worker**     | ECS Fargate Spot (On-demand)     | ECS Fargate Service (상시)   |
| **Queue**         | SQS (wake-up signal)             | SQS (job queue)              |
| **Job Store**     | face_recognition_jobs (Supabase) | face_recognition_jobs (RDS)  |
| **Database**      | Supabase (pgvector 포함)         | RDS PostgreSQL 17 + pgvector |
| **Storage**       | S3                               | 동일                         |
| **CDN**           | CloudFront + OAC                 | 동일                         |

---

## 4. Request Flow

### 이미지 업로드 흐름

```
① POST /media-item/upload-batch        → batch_id 발급
② POST /media-item/presigned-url       → media_item_id + Presigned PUT URL
   └─ API Lambda: media_items (status=01), media_files (role=01) DB 생성
③ 클라이언트: PUT {presigned_url}      → S3 직접 업로드
   └─ S3: original/{familyId}/{albumId}/{date}/{uuid}.jpg
④ S3 PutObject 이벤트 → Resize Lambda
   └─ 리사이즈 → S3 view/, thumbnail/ 업로드
   └─ DB: media_files 생성, media_items.status = '02'
   └─ DB: face_recognition_jobs INSERT (status='01')
   └─ SQS: SendMessage (wake-up signal)
⑤ SQS → trigger-ecs Lambda (Phase 1)
   └─ ECS 실행 중 확인 → 없으면 RunTask
⑥ ECS Fargate Task: Python AI
   └─ face_recognition_jobs에서 배치 fetch (FOR UPDATE SKIP LOCKED)
   └─ S3 view 이미지 다운로드
   └─ InsightFace(buffalo_s) 얼굴 인식 → 512차원 ArcFace 임베딩
   └─ pgvector 코사인 유사도 → identity 매칭 (임계값 0.5)
   └─ face_detections, identity_face_imgs 저장
   └─ 얼굴 크롭(512×512) → S3 identities/ 업로드
   └─ media_items.status = '03' (completed)
   └─ pending job 없으면 Task 자동 종료 (Phase 1)
⑦ 클라이언트: GET /media-item/upload-batch/status (폴링)
```

#### 업로드 상태 천이

```
01 (pending)
    ↓  [Resize Lambda]
02 (processing)
    ↓  [ECS AI]
03 (completed)   ← 정상 완료
04 (failed)      ← 처리 실패 (재처리 대상)
05 (duplicate)   ← 중복 사진 감지
```

### 조회 흐름

```
GET /media-item?group_id=X&from=...&to=...
    └─ album_groups_permissions 권한 확인 (permission='R')
    └─ upload_status = '03' (completed)인 항목만 반환
    └─ media_files JOIN으로 original/thumbnail/view URL 포함
    └─ CloudFront URL로 변환하여 반환

GET /media-item/search?identity_ids=[1,2]&from=...
    └─ face_detections JOIN으로 특정 인물 필터링
    └─ 페이지네이션: 30개씩 (cursor 기반)
```

---

## 5. Multi-Tenancy Strategy

### family_id 기반 설계

```
families
   └─ groups (family 내 그룹)
       └─ users
           └─ albums
               └─ media_items
                   └─ media_files
                   └─ face_detections
                       └─ identities (family 내 인물)
```

### 데이터 격리 전략

| 레이어            | 격리 방법                                                |
| ----------------- | -------------------------------------------------------- |
| **API**           | JWT에서 family_id 추출, 모든 쿼리에 family_id 조건 포함  |
| **앨범 접근**     | `album_groups_permissions` 테이블로 그룹별 R/W 권한 관리 |
| **S3 경로**       | `{familyId}/` prefix로 논리적 분리                       |
| **pgvector 검색** | `WHERE i.family_id = $1` 조건으로 벡터 검색 범위 제한    |

---

## 6. Deployment Topology

### Phase 1 — 네트워크 구성 (VPC 없음)

```
Region: ap-northeast-1 (도쿄)

외부 서비스:
  [ Supabase ]  ← Lambda / ECS가 공용 인터넷으로 접근 (TLS)

AWS:
  [ API Gateway ]           → Go API Lambda (VPC 없음)
  [ S3 ]                    → Resize Lambda (VPC 없음)
  [ SQS ]                   → trigger-ecs Lambda (VPC 없음)
  [ ECS Fargate Spot Task ] → Supabase, S3 접근 (인터넷)
  [ CloudFront + OAC ]      ← 클라이언트 이미지 요청
  [ ECR ]                   → Lambda / ECS Docker 이미지
```

VPC 없음 → Lambda 콜드 스타트 빠름, NAT Gateway 비용 없음

### Phase 2 — 네트워크 구성 (VPC 포함)

```
┌── VPC: 10.0.0.0/16 ──────────────────────────────┐
│  ┌── 퍼블릭 서브넷 ────────────────────────────┐  │
│  │  [ NAT Gateway ]                             │  │
│  └──────────────────────────────────────────────┘  │
│  ┌── 프라이빗 서브넷 ──────────────────────────┐  │
│  │  [ Go API Lambda ]   [ Resize Lambda ]       │  │
│  │  [ ECS Fargate Service (Python AI) ]         │  │
│  │  [ RDS PostgreSQL 17 + pgvector ]            │  │
│  │  VPC Endpoints:                              │  │
│  │    S3 Gateway, ECR API/DKR, CloudWatch       │  │
│  └──────────────────────────────────────────────┘  │
└───────────────────────────────────────────────────┘
```

### 보안 그룹 (Phase 2)

```
sg-api-lambda:   Outbound → sg-rds:5432, S3 VPC Endpoint, 인터넷:443
sg-resize-lambda: Outbound → sg-rds:5432, S3 VPC Endpoint, SQS:443
sg-ecs-ai:       Outbound → sg-rds:5432, S3 VPC Endpoint, SQS:443
sg-rds:          Inbound  ← sg-api-lambda, sg-resize-lambda, sg-ecs-ai (5432)
```

---

## 7. Future Evolution

### Phase 1 → Phase 2 전환 로드맵

| 단계 | 작업                               | 비고                       |
| ---- | ---------------------------------- | -------------------------- |
| 1    | Supabase → RDS 데이터 마이그레이션 | pg_dump / restore          |
| 2    | VPC 구성, RDS 생성                 | Terraform                  |
| 3    | Lambda VPC 연결                    | DATABASE_URL 환경변수 교체 |
| 4    | trigger-ecs Lambda 제거            | ECS Service로 교체         |
| 5    | ECS Service 배포 (SQS 폴링 방식)   | Task → Service 전환        |
| 6    | SQS Auto Scaling 정책 설정         | queue depth 기반           |

### ECS Auto Scaling (Phase 2)

```
SQS queue depth (ApproximateNumberOfMessages)
    → CloudWatch Metric
    → Application Auto Scaling
        → queue depth / 50 = 목표 Task 수
        → 최소 1, 최대 10
```

### 홈서버 이전 시 변경 요소

| AWS 서비스      | 홈서버 대체                             |
| --------------- | --------------------------------------- |
| S3              | MinIO (S3 호환 API)                     |
| Lambda (API)    | Docker 컨테이너                         |
| Lambda (Resize) | 별도 컨테이너 또는 동일 프로세스 고루틴 |
| ECS Fargate     | Docker 컨테이너 (GPU 서버 권장)         |
| RDS PostgreSQL  | 로컬 PostgreSQL + pgvector              |
| SQS             | Redis (List 또는 Stream)                |
| CloudFront      | Nginx reverse proxy                     |

```go
// 교체 가능한 스토리지 인터페이스
type StorageService interface {
    SaveFile(file []byte, filePath string, fileName string) (string, error)
    GeneratePresignedPutURL(key, contentType string, expiresIn time.Duration) (string, error)
}
// STORAGE_TYPE 환경변수로 "s3" / "minio" / "local" 선택
```
