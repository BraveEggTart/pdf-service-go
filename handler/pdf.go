package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"pdf-service-go/pdf"
)

// 默认 PDF 渲染超时（秒）
const defaultPDFTimeoutSec = 30

// Health 返回 GET /health 的 JSON：{"status":"ok"}
func Health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

// PDF 处理 POST /pdf：body 为原始 HTML 或 JSON {"html":"..."}，返回 application/pdf。
func PDF(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	html, err := parseRequestBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if html == "" {
		http.Error(w, "missing or empty html", http.StatusBadRequest)
		return
	}

	timeoutSec := defaultPDFTimeoutSec
	if n := getEnvInt("PDF_TIMEOUT_SECONDS"); n > 0 && n <= 120 {
		timeoutSec = n
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	pdfBytes, err := pdf.RenderHTMLToPDF(ctx, html)
	if err != nil {
		http.Error(w, "pdf render failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pdfBytes)
}

func getEnvInt(key string) int {
	s := os.Getenv(key)
	if s == "" {
		return 0
	}
	n, _ := strconv.Atoi(s)
	return n
}

func parseRequestBody(r *http.Request) (string, error) {
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		var body struct {
			HTML string `json:"html"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			return "", err
		}
		return body.HTML, nil
	}
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}
