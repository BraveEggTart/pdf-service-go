# 构建阶段
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server .

# 运行阶段：Debian + Chromium + 系统字体（支持拉丁/西欧、CJK、阿拉伯、泰语及 Indic/东南亚 等小众语种）
# Railway 默认使用本 Dockerfile，无需内嵌字体即可支持多语言 PDF
FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        chromium \
        fontconfig \
        fonts-noto-core \
        fonts-noto-cjk \
        fonts-noto-cjk-extra \
        fonts-thai-tlwg \
        fonts-deva \
        fonts-beng \
        fonts-guru \
        fonts-gujr \
        fonts-orya \
        fonts-taml \
        fonts-telu \
        fonts-knda \
        fonts-mlym \
        fonts-lao \
        fonts-khmeros \
        fonts-sil-padauk \
        fonts-ddc-uchen \
    && fc-cache -fv \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /server /server
EXPOSE 8080
ENV PORT=8080
ENV CHROME_PATH=/usr/bin/chromium
CMD ["/server"]
