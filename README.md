# PDF Service (Go)

接收 Heroku 主项目发来的 HTML（可含各种小众语言），用 Docker 内 Chromium + 系统字体渲染为 PDF，避免乱码。主项目通过 `PDF_CHROME_URL` 调用本服务。

## 接口

- **POST /pdf** — Body：原始 HTML 或 `{"html":"..."}`；返回 `application/pdf`。
- **GET /health** — 返回 `{"status":"ok"}`。

## 支持语种（Docker 系统字体）

拉丁/西欧、简繁体中文、日韩、阿拉伯、泰语；天城文、孟加拉、Gurmukhi、Gujarati、Oriya、泰米尔、泰卢固、卡纳达、马拉雅拉姆；老挝、藏文、高棉、缅甸等。无需下载或内嵌字体，Docker 镜像内已安装。

## 本地运行

```bash
go mod tidy && go run .
# 或 Docker：docker build -t pdf-service-go . && docker run -p 8080:8080 pdf-service-go
```

默认端口 `8080`，可用环境变量 `PORT`、`PDF_TIMEOUT_SECONDS`（秒）、`CHROME_PATH`。

## Railway 部署

1. 连接 GitHub，Root Directory 设为 `pdf-service-go`。
2. 使用默认 **Dockerfile** 构建。
3. 生成 Public URL；主项目设置 `PDF_CHROME_URL` = `https://<your-app>.railway.app/pdf`。

## 项目结构

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
