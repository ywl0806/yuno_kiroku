# YUNO 프로젝트 개요

## 프로젝트 소개

YUNO는 가족의 추억을 공유하는 사진/동영상 공유 앱입니다. AI 얼굴 인식 기능으로 가족 구성원을 자동으로 분류하고, 가족·그룹 단위 앨범 관리를 지원합니다. Web과 Mobile 앱 모두 지원합니다.

---

## 레포지토리 구성

```
yuno/
├── front/       # Web 프론트엔드 (React + TypeScript)
├── mobile/      # 모바일 앱 (React Native)
├── server/      # 백엔드 API (Go + Echo)
├── ai/          # AI 얼굴인식 서비스 (Python)
├── infra/       # 인프라 (Terraform / AWS)
└── docs/        # 프로젝트 문서
```

---

## 주요 기능

| 기능 | 설명 |
|------|------|
| 사진·동영상 업로드 | Web/Mobile에서 가족 앨범으로 업로드 |
| AI 얼굴인식 | InsightFace(ArcFace)로 얼굴 감지·동일인물 식별 |
| 가족·그룹 관리 | Family 단위 그룹 관리, 초대 링크 발급 |
| 앨범 관리 | 그룹별 접근 권한 앨범 |
| OAuth 로그인 | LINE / Kakao 소셜 로그인 |
| 다국어 지원 | 한국어 / 일본어 (i18n) |
| PWA | 오프라인 지원 (Web 앱) |
| 동영상 처리 | 동영상 업로드 및 재생 지원 |

---

## 기술 스택

| 레이어 | 기술 | 버전 | 용도 |
|--------|------|------|------|
| **Web 프론트엔드** | React | 18.2.0 | UI 프레임워크 |
| | TypeScript | 5.0.2 | 타입 안전성 |
| | Vite + SWC | 7.0.6 | 빌드 도구 |
| | Tailwind CSS | 4.1.11 | 스타일링 |
| | React Router | 6.18.0 | 클라이언트 라우팅 |
| | TanStack Query | 5.7.2 | 서버 상태 관리 |
| | shadcn/ui | - | UI 컴포넌트 |
| | i18next | 25.8.11 | 다국어 |
| **모바일** | React Native | 0.73.5 | 크로스플랫폼 UI |
| | TypeScript | 5.0.4 | 타입 안전성 |
| | React Navigation | 6.1+ | 모바일 라우팅 |
| | NativeWind | 2.0.11 | Tailwind for RN |
| | TanStack Query | 5.24.8 | 서버 상태 관리 |
| **백엔드** | Go | 1.24.0 | 서버 언어 |
| | Echo | 4.11.4 | Web 프레임워크 |
| | sqlc | - | 타입 안전 SQL 코드 생성 |
| | golang-migrate | 4.19.0 | DB 마이그레이션 |
| | govips (libvips) | v2 | 이미지 처리 |
| **데이터베이스** | PostgreSQL | 17 | 메인 DB |
| | pgvector | - | 벡터 유사도 검색 |
| **AI/ML** | InsightFace | 0.7.3 | 얼굴 감지·임베딩 |
| | OpenCV | 4.12.0 | 이미지 처리 |
| **스토리지** | MinIO | latest | S3 호환 오브젝트 스토리지 (로컬 개발) |
| **큐** | ElasticMQ | 1.6.7 | SQS 호환 메시지 큐 (로컬 개발) |

---

## 인증 흐름

```
사용자
  │
  ├─ OAuth (LINE / Kakao)
  │     └─ 콜백 → 서버에서 JWT 발급 (HTTP-only Cookie)
  │
  └─ JWT 토큰
        ├─ Access Token: HTTP-only Cookie
        └─ Refresh Token: DB에 저장
              └─ 401 수신 시 자동으로 /auth/refresh 호출
```

---

## 사진 업로드 흐름

```
1. 클라이언트 → POST /media-item/presigned-url
       ↓ (DB에 media_item 레코드 먼저 생성 후 Presigned URL 반환)
2. 클라이언트 → S3/MinIO에 직접 PUT 업로드
       ↓ (MinIO Webhook 이벤트: original/ prefix)
3. resize-worker 수신 → 리사이즈 (view 2048px, thumbnail 512px)
       ↓ (리사이즈 완료 → SQS face-recognition 큐에 메시지 전송)
4. ai-batch 수신 → InsightFace로 얼굴 감지 → 512차원 임베딩 생성
       ↓ (pgvector로 동일인물 매칭)
5. identity 생성/매핑 → face_detections 저장
```

---

## 로컬 개발 환경

`docker-compose`로 전체 환경을 구성합니다.

```yaml
서비스 구성:
  - app                    : Go 서버            (port: 1323)
  - resize-worker          : 이미지 리사이즈 워커 (port: 1325)
  - ai-batch               : AI 얼굴인식 배치    (SQS 폴링)
  - video-processing-worker: 동영상 처리 워커    (SQS 폴링)
  - postgres               : PostgreSQL 17      (port: 5433)
  - minio                  : S3 호환 스토리지   (port: 9001 / 9090)
  - queue                  : ElasticMQ (SQS)   (port: 9324)
```

자세한 내용은 [development.md](./development.md)를 참조하세요.

---

## 문서 목록

| 파일 | 설명 |
|------|------|
| [overview.md](./overview.md) | 프로젝트 개요 (이 파일) |
| [frontend.md](./frontend.md) | Web 프론트엔드 상세 |
| [mobile.md](./mobile.md) | 모바일 앱 상세 |
| [server.md](./server.md) | 백엔드 서버 상세 |
| [ai.md](./ai.md) | AI 서비스 상세 |
| [development.md](./development.md) | 로컬 개발 환경 설정 |
