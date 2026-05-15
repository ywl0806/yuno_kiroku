# 1. System Architecture Overview

## 1. Introduction

- 프로젝트 목적
- 대상 사용자 (개인 + 가족)
- 설계 목표 (운영 가능성 / 비용 균형)

## 2. High-Level Architecture

- 전체 아키텍처 다이어그램
- 주요 컴포넌트 설명
  - API Server
  - Worker
  - Queue
  - Storage
  - Database

## 3. Request Flow

- 이미지 업로드 흐름 ( 업로드, 리사이즈, 얼굴 인식, 얼굴 크롭, 레코드 생성)
- 조회 흐름

## 4. Multi-Tenancy Strategy

- family_id 기반 설계
- 데이터 격리 전략

⸻

# 2. Upload Pipeline Design

## 1. Design Goals

- 데이터 무결성 보장
- Idempotency
- 재처리 가능 구조
- 비용 효율

## 2. State Machine

- 상태 정의
- 상태 전이 다이어그램
- 전이 조건

## 3. Upload Sequence

- DB 선 생성
- Presigned URL 발급
- 업로드 검증
- SQS 발행

## 4. Idempotency Strategy

- client_upload_id 전략
- hash 기반 보조 검증
- 중복 메시지 처리 방식

## 5. Failure Handling

- 업로드 중단
- S3 실패
- Worker 실패
- 메시지 중복

## 6. Reprocessing Strategy

- failed 상태 재처리
- processing timeout 회수 로직

## 7. Data Consistency Guarantees

- 고아 파일 방지 전략
- 순서 보장 방법

⸻

# 3. AWS Infrastructure Design

...

# 4. Data Model Design

## 1. Design Principles

- 멀티테넌시
- 정합성 우선
- 조회 최적화

## 2. ERD

- 전체 다이어그램

## 3. Core Tables

- families
- users
- albums
- media_items
- media_files
- identities
- face_detections
- upload_batches

(각 테이블별)

- 역할
- 주요 컬럼 설명
- 인덱스 전략

## 4. Vector Search Design

- pgvector 선택 이유
- ivfflat 설정
- embedding dimension
- identity 평균 벡터 전략

## 5. Integrity Constraints

- unique 제약
- foreign key 전략
- soft delete 여부

⸻

# 5. Failure Scenario & Recovery Strategy

## 1. 운영 목표

- 자동 복구 우선
- 데이터 손실 방지

## 2. Failure Categories

- API Layer
- Worker Layer
- Storage
- Database
- Queue

## 3. Scenario Breakdown

### 3.1 Worker Crash

- 원인
- 감지 방법
- 자동 복구
- 수동 조치

### 3.2 Processing Stuck

- timeout 기준
- 회수 알고리즘

### 3.3 S3 업로드 실패

- 재시도 전략
- 클라이언트 UX 처리

### 3.4 DB 장애

- 읽기/쓰기 영향
- 복구 전략

### 3.5 OOM 발생

- 방지 전략
- 격리 전략

## 4. DLQ Strategy

## 5. Disaster Recovery

- 백업 전략
- 복구 절차

⸻

# 6. Cost Analysis & Scaling Strategy

## 1. Current Traffic Assumption

- 일 업로드 용량
- 평균 파일 크기
- 처리 빈도

## 2. AWS Cost Breakdown

- ECS
- RDS
- S3
- SQS
- 데이터 전송

## 3. Monthly Cost Estimation

- 현재 기준
- 5배 증가
- 10배 증가

## 4. Scaling Strategy

- API Scaling
- Worker Scaling
- DB Scaling

## 5. Serverless vs Container 비교

## 6. Home Server Cost Projection

- 전력 비용
- 장비 비용
- 유지관리 비용

⸻

# 7. Security & Multi-Tenancy Model

## 1. Authentication Model

## 2. Authorization Model

## 3. Family Isolation Strategy

## 4. Presigned URL Security

## 5. Data Access Control

## 6. Future RLS 적용 가능성

⸻

# 8. Observability & Monitoring Design

## 1. Logging Strategy

## 2. Metrics Definition

- 처리 시간
- 실패율
- 큐 적체량

## 3. Alert Policy

## 4. 장애 탐지 플로우

## 5. Tracing 전략

⸻
