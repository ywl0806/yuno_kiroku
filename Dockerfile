FROM golang:1.23-alpine

# libde265-dev 설치 (HEIF 관련 라이브러리)
RUN apk update && apk add --no-cache libde265-dev pkgconfig gcc musl-dev g++

ENV CGO_ENABLED=1
WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# air 설정 파일 복사
COPY .air.toml /app/.air.toml

EXPOSE 1323
# air 설정 파일을 사용하여 개발 서버 실행
CMD ["air"]