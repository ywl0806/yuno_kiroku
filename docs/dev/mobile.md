# 모바일 앱

**디렉토리:** `mobile/`

---

## 사용 기술

| 카테고리 | 라이브러리 | 버전 |
|----------|-----------|------|
| 프레임워크 | React Native | 0.73.5 |
| 언어 | TypeScript | 5.0.4 |
| 번들러 | Metro | - |
| 스타일링 | NativeWind (Tailwind for RN) | 2.0.11 |
| 내비게이션 | React Navigation | 6.1+ |
| 서버 상태 관리 | TanStack React Query | 5.24.8 |
| HTTP 클라이언트 | Axios | 1.6.8 |
| 로컬 스토리지 | AsyncStorage | 2.2.0 |
| 미디어 접근 | Camera Roll | 7.4.2 |
| 파일 시스템 | React Native FS | 2.20.0 |
| 동영상 재생 | React Native Video | 5.2.1 |
| UI 아이콘 | React Native Vector Icons | - |
| 사진 목록 | Masonry List | 1.4.2 |

---

## 디렉토리 구조

```
mobile/
├── src/
│   ├── screens/          # 화면 컴포넌트
│   │   ├── HomeScreen    # 타임라인·사진 목록
│   │   ├── LoginScreen   # 로그인
│   │   ├── UploadScreen  # 사진 업로드
│   │   └── SettingsScreen# 설정
│   ├── navigation/       # 내비게이션 설정
│   │   └── 커스텀 탭바
│   ├── services/         # API 서비스
│   │   ├── getPhotos.ts       # 서버에서 사진 조회
│   │   ├── uploadPhoto.ts     # 사진 업로드
│   │   └── getLocalPhotos.ts  # 단말기 내 사진 조회
│   ├── context/          # 전역 상태
│   │   ├── AuthContext   # 인증 상태 관리
│   │   └── TabContext    # 탭 내비게이션 상태
│   ├── components/       # 재사용 가능 컴포넌트
│   ├── hooks/            # 커스텀 훅
│   ├── types/            # TypeScript 타입 정의
│   ├── utils/            # 유틸리티 함수
│   └── lib/              # 헬퍼 함수
├── ios/                  # iOS 네이티브 코드
├── android/              # Android 네이티브 코드
└── package.json
```

---

## 내비게이션 구성

```
RootStack (Native Stack)
├── LoginScreen          # 미인증 시
└── MainTabs (Bottom Tab Navigator)
    ├── Home             # 홈 탭
    ├── Upload           # 업로드 탭
    └── Settings         # 설정 탭
```

- 하단 탭에 커스텀 탭바 사용
- 인증 상태에 따라 루트 전환

---

## 인증

- **OAuth:** LINE / Kakao 소셜 로그인
- **토큰 저장:** `AsyncStorage` (영속화)
- **전역 상태:** `AuthContext`로 관리
- **Axios 인터셉터:** 요청 시 JWT 자동 부여

---

## 상태 관리

### 서버 상태: TanStack React Query

- 사진 목록 등 API 데이터를 캐시·관리

### 클라이언트 상태: Context API

- `AuthContext`: 인증 토큰·유저 정보
- `TabContext`: 탭 내비게이션 상태

---

## 미디어 기능

### 단말기 내 사진 조회

- `Camera Roll`로 단말기 카메라롤에 접근
- 사진 메타데이터 (촬영일시·GPS 정보) 조회

### 사진 업로드

1. 카메라롤에서 사진 선택
2. 백그라운드에서 서버로 업로드
3. React Query로 캐시 무효화 후 화면 갱신

### 동영상 재생

- `React Native Video`로 서버의 동영상 재생

---

## 개발 명령어

```bash
# 의존성 설치
npm install

# iOS Pod 설치
cd ios && pod install

# Metro 개발 서버 시작
npm start

# iOS 실행
npm run ios

# Android 실행
npm run android

# 린트
npm run lint

# 테스트
npm run test
```

---

## 지원 플랫폼

| 플랫폼 | 지원 |
|--------|------|
| iOS | ✅ |
| Android | ✅ |
