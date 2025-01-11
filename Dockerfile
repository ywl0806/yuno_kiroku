FROM golang:1.23-alpine

# libde265-dev 설치 (HEIF 관련 라이브러리)
RUN apk update && apk add --no-cache libde265-dev pkgconfig

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]