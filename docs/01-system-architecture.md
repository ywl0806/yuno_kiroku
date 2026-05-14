# 📘 System Architecture Overview

## 1. Introduction

### 프로젝트 목적

YUNO는 가족 단위의 사진·동영상을 공유하고, 얼굴 인식 기반으로 자동 분류하는 미디어 관리 플랫폼이다.
단순한 앨범 공유를 넘어, 가족 구성원의 얼굴을 자동으로 인식하고 인물별 검색 및 필터링을 제공한다.
동영상 업로드 및 재인코딩도 지원한다.

### 대상 사용자

- **가족 단위**: 부모/자녀가 같은 Family에 속해 사진·동영상을 공유
- **멀티 디바이스**: Web, iOS, Android에서 동일한 경험 제공
- **비기술 사용자**: 업로드 후 자동 처리, 별도 조작 없이 얼굴 분류 완료

### 설계 목표

| 목표                       | 내용                                                   |
| -------------------------- | ------------------------------------------------------ |
| **비용 최소화**            | 유휴 시간 비용 최소화, Serverless 기반 설계            |
| **단계적 확장**            | 개인 운용 → 실 서비스 전환 시 컴포넌트 교체만으로 대응 |
| **클라우드 종속성 최소화** | 인터페이스 추상화로 AWS → 홈서버 이전 가능한 구조      |
| **데이터 무결성**          | 미디어 처리 실패 시 재처리 가능, 고아 파일 방지        |
| **역할 분리**              | ML 추론(Python)과 DB/스토리지 처리(Go) 명확히 분리     |

---

## 2. 아키텍처 단계

### Phase 1 — 개인 운용 (현재)

**목표**: 소규모 사용, 비용 최소화.

```
[ 클라이언트 ]
      │ HTTPS
      ▼
[ API Gateway + Go API Lambda ]
      │  DB: Supabase PostgreSQL (외부)
      │  S3: Presigned URL 발급
      ▼
[ S3: yuno-media-bucket ]
  original/ prefix → S3 이벤트 알림
      │
      ▼
[ Resize Worker (Go + govips) ]  ← HTTP 서버, S3/MinIO Webhook 수신
  이미지:
    - 리사이즈: view(2048px), thumbnail(512px)
    - DB: media_items status=completed
    - SQS: face-recognition 큐에 job 발행
  동영상:
    - SQS: video-processing 큐에 job 발행
      │
      ├─▶ [ SQS: face-recognition-queue ]
      │         ↓
      │   [ AI Batch (Python + InsightFace) ]
      │     - SQS 폴링
      │     - 얼굴 감지 + 512차원 임베딩 계산
      │     - subprocess: face-recognition-worker (Go CLI) 호출
      │         ↓
      │   [ face-recognition-worker (Go CLI) ]
      │     - pgvector 코사인 유사도 → identity 매칭/생성
      │     - face_detections INSERT
      │     - 얼굴 크롭(512×512) → S3 identities/ 저장
      │     - identity_face_imgs INSERT
      │
      └─▶ [ SQS: video-processing-queue ]
                ↓
          [ Video Processing Worker (Go + ffmpeg) ]
            - SQS 폴링
            - 비디오 리사이즈 (최대 1280px, H.264)
            - 썸네일 추출 (1초, WebP)
            - DB: media_files 저장
            - DB: media_items status=completed
      │
      ▼
[ Supabase PostgreSQL + pgvector ]
```

---

### Phase 2 — 실 서비스 (미래)

**목표**: 다수 사용자 대응, 안정성 우선.

```
[ 클라이언트 ]
      │ HTTPS
      ▼
[ API Gateway + Go API Lambda (VPC) ]
      │  RDS PostgreSQL (VPC 내부)
      │  S3: Presigned URL 발급
      ▼
[ S3: yuno-media-bucket ]
      │ S3 이벤트 → Resize Worker (ECS Service)
      ▼
[ Resize Worker ECS Service (Go + govips) ]
  이미지/동영상 분기 후 SQS 발행
      │
      ├─▶ [ SQS: face-recognition-queue ]
      │         ↓ ECS AI Service (상시 가동, Auto Scaling)
      │   [ AI Batch + face-recognition-worker ]
      │
      └─▶ [ SQS: video-processing-queue ]
                ↓ ECS Video Service (상시 가동)
          [ Video Processing Worker ]
      │
      ▼
[ RDS PostgreSQL ]

```

---

## 3. 주요 컴포넌트

| 컴포넌트                    | 기술                                    | 역할                                                                    |
| --------------------------- | --------------------------------------- | ----------------------------------------------------------------------- |
| **API Server**              | Go + Echo (Lambda 어댑터)               | REST API, Presigned URL 발급                                            |
| **Resize Worker**           | Go + govips (Lambda function)           | S3 PutObject Event 수신(로컬: MinIO Webhook), 이미지 리사이즈, SQS 발행 |
| **AI Batch**                | Python + InsightFace                    | SQS 폴링, 얼굴 감지·임베딩 계산, Go CLI 호출                            |
| **face-recognition-worker** | Go CLI 바이너리                         | identity 매칭, face_detections INSERT, 얼굴 크롭                        |
| **Video Worker**            | Go + ffmpeg                             | SQS 폴링, 동영상 리사이즈·썸네일 추출                                   |
| **Queue**                   | SQS (로컬: ElasticMQ)                   | face-recognition / video-processing 큐                                  |
| **Database**                | Supabase / RDS PostgreSQL 17 + pgvector | 메타데이터, 벡터 검색                                                   |
| **Storage**                 | S3 (로컬: MinIO)                        | 미디어 파일 저장                                                        |
| **CDN**                     | CloudFront + OAC                        | 미디어 파일 서빙                                                        |

---

## 4. Request Flow

### 이미지 업로드 흐름

```
① POST /media-item/presigned-url
   └─ API: media_items(status=01) + media_files(role=original) DB 생성
   └─ API: Presigned PUT URL 발급
② 클라이언트: PUT {presigned_url} → S3 직접 업로드
   └─ S3: original/{familyId}/{mediaItemId}.{ext}
③ S3 PutObject Event → resize lambda function (로컬: MinIO Webhook)
   └─ media_items.status = '02' (processing)
   └─ 원본 다운로드 → EXIF 파싱 (촬영일시·GPS)
   └─ view(2048px) + thumbnail(512px) 리사이즈 → S3 업로드
   └─ media_files 생성 (view, thumbnail)
   └─ media_items.status = '03' (completed) ← 이 시점부터 사진 조회 가능
   └─ SQS face-recognition 큐에 발행 {media_item_id, family_id, view_storage_key}
④ ai-batch (Python) — SQS 폴링
   └─ view 이미지 S3 다운로드 (임시 파일)
   └─ InsightFace(buffalo_s) 얼굴 감지 → 512차원 ArcFace 임베딩
   └─ subprocess: face-recognition-worker (Go CLI)
      stdin: {media_item_id, family_id, view_image_path, faces}
⑤ face-recognition-worker (Go CLI)
   └─ for each face:
      └─ advisory lock (family 단위)
      └─ pgvector 코사인 유사도 → identity 매칭/신규 생성 (임계값 0.5)
      └─ face_detections INSERT (위치 + 임베딩)
      └─ view 이미지에서 얼굴 크롭(512×512) → S3 identities/ 저장
      └─ identity_face_imgs INSERT
```

### 동영상 업로드 흐름

```
① ~ ② 이미지와 동일 (Presigned URL → S3 업로드)
③ S3 PutObject Event → resize lambda function (로컬: MinIO Webhook)
   └─ 동영상 확장자 감지 (mp4, mov, avi, mkv, webm, m4v)
   └─ media_items.status = '02' (processing)
   └─ SQS video-processing 큐에 발행 {media_item_id, family_id, original_storage_key, ...}
④ video-processing-worker (Go) — SQS 폴링
   └─ 원본 동영상 S3 다운로드
   └─ ffmpeg 리사이즈 (최대 1280px, H.264 CRF23, faststart)
   └─ ffmpeg 썸네일 추출 (1초 지점, WebP)
   └─ thumbnail + video S3 업로드
   └─ media_files 생성 (thumbnail + video)
   └─ media_items.status = '03' (completed)
```

### 업로드 상태 천이

```
01 (pending)
    ↓  [resize-worker: 처리 시작]
02 (processing)
    ↓  [resize-worker: 이미지 완료 / video-worker: 동영상 완료]
03 (completed)   ← 정상 완료, 미디어 조회 가능
04 (failed)      ← 처리 실패
```

### 조회 흐름

```
GET /media-item?from=...&to=...
    └─ upload_status = '03' (completed)인 항목만 반환
    └─ media_files JOIN으로 thumbnail/view/video URL 포함

GET /media-item/search?identity_ids=[1,2]
    └─ face_detections JOIN으로 특정 인물 필터링
    └─ 페이지네이션: cursor 기반
```

---

## 5. AI 처리 아키텍처 — Python + Go 하이브리드

Python과 Go의 역할을 명확히 분리한다.

```
[ ai-batch (Python) ]          [ face-recognition-worker (Go CLI) ]
  ML 추론만 담당                  DB/스토리지 처리 담당
  ─────────────────              ────────────────────────────────
  InsightFace 모델 로딩           pgvector 코사인 유사도 검색
  얼굴 감지 (bbox)                identity 매칭 / 신규 생성
  512차원 임베딩 계산             face_detections INSERT
  임시 파일로 저장                얼굴 크롭 (512×512)
         │                       S3 업로드 (identities/)
         │ subprocess stdin JSON  identity_face_imgs INSERT
         └─────────────────────▶
```

이 분리의 장점:

- Python은 무거운 ML 라이브러리(InsightFace, OpenCV)에 집중
- Go는 타입 안전한 DB 처리, pgvector, 스토리지 작업 담당
- 각 컴포넌트를 독립적으로 배포·교체 가능

---

## 6. Multi-Tenancy Strategy

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

| 레이어            | 격리 방법                                                   |
| ----------------- | ----------------------------------------------------------- |
| **API**           | JWT에서 family_id 추출, 모든 쿼리에 family_id 조건 포함     |
| **앨범 접근**     | `album_groups_permissions` 테이블로 그룹별 R/W 권한 관리    |
| **S3 경로**       | `{familyId}/` prefix로 논리적 분리                          |
| **pgvector 검색** | `WHERE family_id = $1` 조건으로 벡터 검색 범위 제한         |
| **Advisory Lock** | `family_id` 해시 기반 잠금으로 동시 identity 생성 충돌 방지 |

---

---
