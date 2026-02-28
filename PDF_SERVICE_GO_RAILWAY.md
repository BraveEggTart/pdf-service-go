# Go PDF 服务技术文档（Railway 部署）

Heroku 主项目发送 HTML 到本服务，HTML 中可能包含各种小众语言；本服务用 Docker 内 Chromium + 系统字体渲染为 PDF，不依赖内嵌字体，避免乱码。

---

## 1. 服务概述

| 项目 | 说明 |
|------|------|
| **功能** | 接收 HTML，用 Headless Chromium 渲染为 PDF，返回二进制。支持复杂脚本（印地语、阿拉伯语、泰语、天城文、孟加拉、泰米尔、藏文、高棉、缅甸等）。 |
| **调用方** | Heroku 主项目，环境变量 `PDF_CHROME_URL` 指向本服务 `/pdf`。 |
| **部署** | Railway，Dockerfile（Debian + Chromium + 系统字体）。 |

---

## 2. 支持语种（Docker 系统字体）

镜像内通过 apt 安装字体包，Chromium 自动使用，无需下载 TTF 或内嵌。

| 类别 | 语种/脚本 |
|------|-----------|
| 基础 | 拉丁/西欧、阿拉伯、泰语 |
| CJK | 简体/繁体中文、日、韩 |
| Indic | 天城文、孟加拉、Gurmukhi、Gujarati、Oriya、泰米尔、泰卢固、卡纳达、马拉雅拉姆 |
| 东南亚等 | 老挝、藏文、高棉、缅甸、Arabic Supplement |

---

## 3. API

- **GET /health** — 200，`{"status":"ok"}`。
- **POST /pdf** — Body：`text/html` 或 JSON `{"html":"..."}`；成功 200，`Content-Type: application/pdf`；错误 400/500。

---

## 4. 项目结构

```
pdf-service-go/
├── go.mod
├── main.go
├── handler/pdf.go
├── pdf/render.go
├── Dockerfile
├── .dockerignore
└── README.md
```

---

## 5. Docker 与 Railway

- 使用仓库内 **Dockerfile**（Debian bookworm-slim + chromium + 多语种字体包）。
- Railway：Deploy from GitHub，Root Directory = `pdf-service-go`，构建后生成 Public URL；主项目设置 `PDF_CHROME_URL` = `https://<your-app>.railway.app/pdf`。

---

## 6. 环境变量

| 变量 | 说明 |
|------|------|
| `PORT` | 监听端口，Railway 注入，默认 8080。 |
| `PDF_TIMEOUT_SECONDS` | PDF 渲染超时（秒），默认 30。 |
| `CHROME_PATH` | 可选，Chromium 路径；Dockerfile 已设 `/usr/bin/chromium`。 |

---

## 7. 参考

- [chromedp](https://github.com/chromedp/chromedp)
- [Railway Docker](https://docs.railway.app/deploy/dockerfiles)
