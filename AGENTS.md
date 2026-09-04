# AGENTS.md

This file provides guidance to Codex (Codex.ai/code) when working with code in this repository.

## 프로젝트 개요

YUNO는 가족 사진 공유 플랫폼으로 다음 컴포넌트로 구성됩니다:
- `front/` - React + TypeScript + Vite (Web, 포트 5155)
- `server/` - Go + Echo (Backend API, 포트 1323)
- `mobile/` - React Native (iOS/Android)
- `ai/` - Python FastAPI (얼굴인식 서비스, 포트 8000)
- `infra/` - Terraform AWS 인프라

## 개발 명령어

### Frontend (`front/`)
```bash
npm run dev       # 개발 서버 (localhost:5155)
npm run build     # 프로덕션 빌드
npm run lint      # ESLint
npm run preview   # 빌드 결과 미리보기
```

### Backend (`server/`)
```bash
# Docker 서비스 시작 (PostgreSQL, MinIO, AI 서비스 포함)
docker compose up -d
docker compose logs -f app ai   # 로그 확인

# DB 작업
make migrate      # 마이그레이션 실행
make seed         # 시드 데이터 삽입
make destroy      # DB 초기화
make refresh      # destroy + migrate + seed

# 코드 생성
make sqlc         # sqlc로 DB 쿼리 코드 생성
make swag         # Swagger 문서 생성

# 포트 정리
make kill         # 1323 포트 프로세스 종료
```

### Mobile (`mobile/`)
```bash
npm start              # Metro 개발 서버
npm run ios            # iOS 시뮬레이터
npm run android        # Android
npm run lint
npm run test           # Jest
```

## 아키텍처

### Frontend 구조
- **`page/`** - 라우트에 대응하는 페이지 컴포넌트 (thin wrapper)
- **`feature/`** - 도메인별 컨테이너/훅/타입/상수 모음 (비즈니스 로직)
- **`components/`** - 재사용 가능한 UI 컴포넌트 (`ui/`는 shadcn/ui)
- **`service/`** - API 호출 함수
- **`lib/`** - Axios 인스턴스, React Query 클라이언트, 세션스토리지 유틸

**라우팅**: `front/src/router.tsx` — `/login`과 `/login/callback`을 제외한 모든 경로는 `DefaultLayout`(인증 필수)으로 감쌈

**API 호출**: 인증 필요 시 `MyAxiosWithAuth`, 불필요 시 `MyAxios` 사용 (`front/src/lib/my-axios.ts`). 401 응답 시 자동으로 토큰 제거 후 `/login`으로 리다이렉트.

**상태관리**: TanStack Query로 서버 상태 관리. 커스텀 훅은 `feature/<domain>/hooks/` 또는 `front/src/hooks/`에 위치.

### Backend 구조 (Go)
Handler → Service → Store → DB (sqlc 생성 쿼리)

- **`internal/handlers/`** - HTTP 요청/응답 처리
- **`internal/services/`** - 비즈니스 로직
- **`internal/store/`** - 데이터 접근 레이어
- **`internal/db/`** - sqlc가 생성한 타입 안전 쿼리 코드 (직접 수정 금지)

Swagger UI: `http://localhost:1323/swagger/index.html`

### 인증 흐름
1. LINE/Kakao OAuth → 서버에서 JWT 발급
2. 콜백 URL로 토큰 전달 (`/login/callback?token=...`)
3. `setAuthToken(token)` 호출 → `localStorage`에 저장 + Axios 헤더 설정

### 다국어 (i18n)
- `front/src/i18n/locales/kr.json` (한국어), `jp.json` (일본어)
- `useTranslation()` 훅으로 사용
- 언어 코드: `ko`/`kr`/`ko-KR` → 한국어, `ja`/`jp`/`ja-JP` → 일본어

## 주요 설정

### 환경변수
- Frontend: `front/.env` — `VITE_API_URL="/api"` (Vite proxy를 통해 1323으로 포워딩)
- Backend: `server/.env` — `.env.example` 참고

### Vite Proxy
`/api` → `http://127.0.0.1:1323`, `/uploads` → `http://127.0.0.1:1323`

### 코드 스타일 (Frontend)
- Prettier: `printWidth: 120`, `singleQuote: true`, `semi: false`
- Tailwind class 정렬: `prettier-plugin-tailwindcss`
- Import 정렬: `prettier-plugin-organize-imports`

## 핵심 파일 참조

| 목적 | 파일 |
|------|------|
| 라우팅 정의 | `front/src/router.tsx` |
| API 엔드포인트 상수 | `front/src/consts/api-route.ts` |
| 전역 타입 정의 | `front/src/types/index.ts` |
| Axios 인스턴스 | `front/src/lib/my-axios.ts` |
| React Query 클라이언트 | `front/src/lib/query-client.ts` |
| 인증 토큰 저장 | `front/src/feature/auth/lib/set-auth-token.ts` |
| 번역 (한국어) | `front/src/i18n/locales/kr.json` |
| 번역 (일본어) | `front/src/i18n/locales/jp.json` |
| 색상 정의 | `front/src/colors.ts` |

## 데이터베이스

PostgreSQL 17 + pgvector (포트 5433)

주요 테이블: `families`, `groups`, `users`, `albums`, `media_items`, `media_files`, `identities`, `face_detections` (512차원 벡터), `invite_tokens`, `refresh_tokens`

DB 스키마 변경 시 반드시 `make sqlc`로 코드 재생성 필요.

## 스토리지

- 개발: MinIO (S3 호환, 포트 9000/9090)
- 프로덕션: AWS S3
- `server/.env`의 `STORAGE_TYPE` 값으로 전환 (`minio` / `s3` / `local`)
