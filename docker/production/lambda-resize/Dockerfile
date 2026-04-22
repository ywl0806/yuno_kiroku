# ---- Build Stage ----
FROM golang:1.24-bullseye AS builder

# libvips 빌드 의존성 설치 (debian bullseye = glibc 2.31, al2023 = glibc 2.34 이므로 호환됨)
RUN apt-get update && apt-get install -y --no-install-recommends \
    libvips-dev \
    libheif-dev \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

ENV CGO_ENABLED=1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -tags lambda.norpc -ldflags="-s -w" -o bootstrap ./cmd/lambda/resize/

# ---- Runtime Stage ----
FROM public.ecr.aws/lambda/provided:al2023

COPY --from=builder /usr/lib/aarch64-linux-gnu/libvips.so* /usr/lib64/
COPY --from=builder /usr/lib/aarch64-linux-gnu/libheif.so* /usr/lib64/
COPY --from=builder /usr/lib/aarch64-linux-gnu/libglib-2.0.so* /usr/lib64/
COPY --from=builder /usr/lib/aarch64-linux-gnu/libgobject-2.0.so* /usr/lib64/
COPY --from=builder /usr/lib/aarch64-linux-gnu/libexpat.so* /usr/lib64/
COPY --from=builder /app/bootstrap /var/runtime/bootstrap

CMD ["bootstrap"]
