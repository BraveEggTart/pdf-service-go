package main

import (
	"log"
	"net/http"
	"os"

	"pdf-service-go/handler"

	"github.com/go-chi/chi/v5"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := chi.NewRouter()
	r.Use(withCORS)
	r.Get("/health", handler.Health)
	r.Post("/pdf", handler.PDF)

	addr := ":" + port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

// withCORS 为所有请求添加简单的 CORS 头，便于从浏览器/本地 HTML 调用生产接口。
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// 预检请求直接返回
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
