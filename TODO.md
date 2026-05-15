# YUNO TODO

## 1. DLQ 적재 시 알림 발송

SQS DLQ(face-recognition-dlq, video-processing-dlq)에 메시지가 쌓이면 운영자에게 즉시 알림을 보낸다.

- [ ] CloudWatch Alarm 추가
  - `ApproximateNumberOfMessagesVisible` > 0 조건으로 face-recognition-dlq / video-processing-dlq 각각 알람 생성
  - `infra/modules/cloudwatch/main.tf`에 추가
- [ ] SNS Topic 생성 및 CloudWatch Alarm Action 연결
  - `infra/modules/sns/` 모듈 생성
  - CloudWatch Alarm → SNS Topic → Email Subscription
- [ ] 알림 이메일 수신자 설정
  - SSM Parameter Store에 이메일 주소 저장, Terraform으로 참조
- [ ] 알림 메시지 포맷 정의
  - 큐 이름, 적재 메시지 수, 발생 시각 포함
- [ ] Terraform 적용 및 테스트
  - DLQ에 수동으로 메시지 투입 후 이메일 수신 확인

---

## 2. DB 백업 전략 수립

Supabase PostgreSQL 데이터 보호 및 복구 절차를 체계화한다.

- [ ] Supabase 백업 설정 확인 및 문서화
  - Free 플랜: 일별 백업 7일 보존 여부 확인
  - Pro 플랜 전환 시 PITR(Point-in-Time Recovery) 활성화 검토
- [ ] 백업 보존 기간 정책 결정
  - 일별 스냅샷 보존 기간 (권장: 30일)
  - 월별 장기 보존 여부
- [ ] 복구 절차 문서화
  - PITR 기반 복구 단계별 절차 작성 (`docs/05-failure-scenario.md` 보완)
  - 복구 목표 수치 확정: RPO / RTO
- [ ] 복구 훈련 계획
  - 분기 1회 복구 드릴 일정 수립
  - 복구 완료 기준 정의 (데이터 무결성 검증 쿼리)
- [ ] S3 미디어 파일 보호 전략
  - S3 버저닝 활성화 여부 결정 (비용 vs 안전성 트레이드오프)
  - 원본 파일 수명주기 정책: 1년 이상 경과 시 Glacier 이전 검토
- [ ] 백업 모니터링
  - Supabase 백업 실패 알림 수신 설정 확인

---

## 3. 로깅 전략 수립

로그의 목적을 먼저 정의하고, 포맷·지표·보존 정책을 일관되게 정립한다.

### 3-1. 로깅 목적 정의

- [ ] 컴포넌트별 로깅 목적 명확화
  - API Server: 요청 추적, 인증 실패, 비즈니스 오류
  - Resize Worker: 처리 성공/실패, 소요 시간
  - AI Batch: 얼굴 감지 건수, 임베딩 처리량, subprocess 오류
  - Video Worker: 트랜스코딩 성공/실패, 파일 크기
- [ ] 로그 레벨 정책 확정 (INFO / WARN / ERROR 기준)

### 3-2. 로그 포맷 표준화

- [ ] 구조화 로그 포맷 결정 (JSON 권장)
  - 공통 필드: `timestamp`, `level`, `component`, `media_item_id`, `family_id`, `message`
- [ ] Go 로거 라이브러리 선택 및 적용
  - `slog` (표준 라이브러리) 또는 `zerolog` 검토
  - API Server, Resize Worker, Video Worker, face-recognition-worker 일괄 적용
- [ ] Python (AI Batch) 로그 포맷 통일
  - `structlog` 또는 `logging` + JSON formatter 적용
- [ ] 민감 정보 마스킹 규칙 정의
  - JWT 토큰, OAuth 코드, DB 연결 문자열 로그 출력 금지

### 3-3. 운영 지표 정의

- [ ] 업로드 파이프라인 커스텀 지표 구현
  - 이미지/동영상 처리 성공률
  - 단계별 평균 소요 시간 (presigned → completed)
  - 일별 실패 건수 (`upload_status = 04`)
- [ ] CloudWatch Custom Metrics 발행 (EMF 또는 PutMetricData)
- [ ] CloudWatch 대시보드 생성
  - SQS 적체량, ECS 태스크 수, 처리 성공률 한 화면에 표시

### 3-4. 로그 보존 전략

- [ ] 컴포넌트별 보존 기간 확정 (현재 30일 일괄 → 용도별 세분화 검토)
  - 오류 로그: 90일
  - 정상 처리 로그: 30일
- [ ] 장기 보존이 필요한 로그 S3 아카이빙 설정 (CloudWatch Logs → S3 Export)
- [ ] 로그 쿼리 가이드 작성 (CloudWatch Logs Insights 쿼리 예제)

---

## 4. 장애 대응 전략 수립

장애 탐지 → 알림 → 대응 → 복구의 전체 흐름을 체계화한다.

### 4-1. 알림 체계 구축

- [ ] 알림 채널 결정 (이메일)
- [ ] 심각도별 알림 대상 정의
  - Critical: 즉시 호출 (DLQ 적재, API 오류율 급증)
  - Warning: 이메일 (처리 지연, 메모리 사용률 높음)
- [ ] 추가 CloudWatch Alarm 구현
  - API Lambda Errors > 50건/5분
  - SQS `ApproximateAgeOfOldestMessage` > 30분
  - ECS Task MemoryUtilization > 80%
  - DLQ 메시지 수 > 0 (1번 항목과 연계)

### 4-2. 장애 탐지 자동화

- [ ] 처리 stuck 감지 자동화
  - `upload_status = 02` 상태가 1시간 이상 지속되는 항목 주기적 탐지
  - Lambda scheduled event 또는 Supabase pg_cron으로 구현
- [ ] Orphan 레코드 정리 자동화
  - `upload_status = 01`이 24시간 이상 경과한 레코드 자동 삭제
  - 삭제 전 S3 파일 존재 여부 확인
- [ ] Resize Worker Webhook 누락 감지
  - `upload_status = 01`인데 S3에 파일이 존재하는 경우 탐지 및 재처리 트리거

### 4-3. 대응 절차 문서화

- [ ] Runbook 작성 (장애 유형별 대응 가이드)
  - DLQ 메시지 재처리 절차
  - Stuck 미디어 수동 복구 절차
  - Resize Worker 수동 재처리 (`POST /resize`) 절차
- [ ] 장애 등급 정의
  - P1: 전체 서비스 불가 (DB 장애, Lambda 다운)
  - P2: 신규 업로드 처리 불가 (Worker 장애)
  - P3: 일부 미디어 처리 지연 (DLQ 적재)

---

## 추가 항목

### 5. 원본 파일 수명주기 정책

- [ ] S3 Lifecycle Rule 구현
  - `original/` prefix: 1년 경과 시 Glacier Instant Retrieval로 이전
  - 비용 절감 효과 계산 (S3 Standard vs Glacier 단가 비교)
- [ ] `infra/modules/s3/main.tf`에 lifecycle 블록 추가
