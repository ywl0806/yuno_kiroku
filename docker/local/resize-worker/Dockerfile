FROM golang:1.24-alpine

# libvips 및 빌드 도구 설치
# community 저장소 활성화 (libvips는 community 저장소에 있음)
RUN apk add --no-cache \
    vips-dev \
    vips \
    libheif-dev \
    libheif \
    pkgconfig \
    gcc \
    musl-dev \
    g++ \
    glib-dev \
    expat-dev

ENV CGO_ENABLED=1

WORKDIR /app

RUN go install github.com/air-verse/air@v1.61.5

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 1325

# air 설정 파일 복사
COPY .air.resize.toml /app/.air.resize.toml

# air 설정 파일을 사용하여 개발 서버 실행
CMD ["air", "-c", ".air.resize.toml"]