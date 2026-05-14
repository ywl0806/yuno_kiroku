# Web 프론트엔드

**디렉토리:** `front/`

---

## 사용 기술

| 카테고리 | 라이브러리 | 버전 |
|----------|-----------|------|
| UI 프레임워크 | React | 18.2.0 |
| 언어 | TypeScript | 5.0.2 |
| 빌드 도구 | Vite + SWC | 7.0.6 |
| 스타일링 | Tailwind CSS | 4.1.11 |
| UI 컴포넌트 | shadcn/ui (Radix UI) | - |
| 라우팅 | React Router DOM | 6.18.0 |
| 서버 상태 관리 | TanStack React Query | 5.7.2 |
| 폼 관리 | React Hook Form | 7.61.1 |
| 유효성 검사 | Zod | 4.0.10 |
| HTTP 클라이언트 | Axios | 1.6.0 |
| 다국어 | i18next + react-i18next | 25.8.11 |
| 아이콘 | Lucide React, MUI Icons | - |
| 캐러셀 | Swiper | 11.1.0 |
| 포토 레이아웃 | react-photo-album | 2.3.1 |
| PWA | vite-plugin-pwa | 1.0.2 |

---

## 디렉토리 구조

```
front/
├── src/
│   ├── feature/          # 도메인별 컨테이너·훅·타입·상수
│   │   ├── auth/         # 인증
│   │   ├── home/         # 홈 (타임라인)
│   │   ├── upload/       # 사진 업로드
│   │   ├── settings/     # 설정
│   │   ├── search/       # 검색
│   │   ├── recent/       # 최근 업로드
│   │   ├── like/         # 좋아요
│   │   └── tag/          # 태그
│   ├── page/             # 라우트에 대응하는 페이지 컴포넌트 (thin wrapper)
│   ├── components/       # 재사용 가능한 UI 컴포넌트
│   │   └── ui/           # shadcn/ui 컴포넌트
│   ├── service/          # API 호출 함수
│   ├── lib/              # 유틸리티
│   │   ├── my-axios.ts   # Axios 인스턴스 (인증 토큰 자동 처리)
│   │   └── query-client.ts # React Query 클라이언트 설정
│   ├── i18n/             # 국제화 리소스
│   │   └── locales/
│   │       ├── kr.json   # 한국어
│   │       └── jp.json   # 일본어
│   ├── hooks/            # 공통 커스텀 훅
│   ├── types/            # TypeScript 타입 정의
│   ├── consts/           # API 라우트 등 상수
│   ├── providers/        # Context 프로바이더
│   ├── router.tsx        # 라우팅 설정
│   └── App.tsx           # 루트 컴포넌트
├── public/               # 정적 파일
├── package.json
├── tsconfig.json
└── vite.config.ts
```

---

## 라우팅 구성

`/login`과 `/login/callback`을 제외한 모든 경로는 `DefaultLayout`(인증 필수)으로 감쌉니다.

```
/login                            # 로그인 페이지
/login/callback                   # OAuth 콜백
/                                 # DefaultLayout (인증 필수)
├── /                             # 홈 (사진 타임라인)
├── /:date                        # 날짜별 사진
├── /recent                       # 최근 업로드 배치 목록
├── /recent/:batchId              # 배치 상세
├── /search                       # 검색
├── /upload                       # 사진 업로드
├── /logout                       # 로그아웃
└── /settings                     # 설정
    ├── /settings/app             # 앱 설정
    ├── /settings/account         # 계정 설정
    ├── /settings/family/new      # 패밀리 생성 (Admin)
    ├── /settings/family/:id/edit # 패밀리 편집 (Admin)
    ├── /settings/family/:id/invite  # 패밀리 초대 (Admin)
    ├── /settings/album/new       # 앨범 생성 (Admin)
    ├── /settings/album/:id/edit  # 앨범 편집 (Admin)
    ├── /settings/member/invite   # 멤버 초대 (Admin)
    ├── /settings/member/:id/edit # 멤버 편집 (Admin)
    ├── /settings/kid/new         # 아이 추가 (Admin)
    └── /settings/kid/:id/edit    # 아이 편집 (Admin)
```

---

## 인증

- **OAuth:** LINE / Kakao 소셜 로그인
- **토큰 관리:** HTTP-only Cookie (서버에서 Set-Cookie)
- **Axios 설정:** `withCredentials: true`로 쿠키 자동 전송
- **토큰 갱신:** 401 응답 시 `/auth/refresh`로 자동 갱신 후 재요청
- **갱신 실패:** `/login`으로 리다이렉트

```typescript
// 인증 필요: MyAxiosWithAuth (withCredentials: true)
// 인증 불필요: MyAxios
import { MyAxiosWithAuth, MyAxios } from '@/lib/my-axios'
```

---

## API 통신

- **베이스 URL:** `VITE_API_URL` 환경변수 (`/api`)
- **Vite Proxy:** `/api` → `http://127.0.0.1:1323`
- **Echo 배열 쿼리 직렬화:** `key=1&key=2` 형식 (Axios 기본 `key[]=1` 대신)

---

## 상태 관리

### 서버 상태: TanStack React Query

- API 데이터 fetching·캐시·동기화
- 커스텀 훅으로 쿼리 캡슐화 (`feature/<domain>/hooks/`)

### 클라이언트 상태: Context API

- 업로드 처리 상태: `UploadPhotoProvider`

---

## 다국어 (i18n)

- 지원 언어: 한국어(`kr`) / 일본어(`jp`)
- `useTranslation()` 훅 사용
- 번역 파일: `src/i18n/locales/kr.json`, `jp.json`

---

## PWA (Progressive Web App)

- `vite-plugin-pwa`로 서비스 워커 생성
- 오프라인 캐시 지원
- 홈 화면 추가 가능

---

## 개발 명령어

```bash
npm run dev       # 개발 서버 (localhost:5155)
npm run build     # 프로덕션 빌드
npm run lint      # ESLint
npm run preview   # 빌드 결과 미리보기
```
