FROM golang:1.24-alpine


WORKDIR /app

RUN go install github.com/air-verse/air@v1.61.5

# sqlc 설치
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# air 설정 파일 복사
COPY .air.toml /app/.air.toml

EXPOSE 1323
# air 설정 파일을 사용하여 개발 서버 실행
CMD ["air", "-c", ".air.toml"]