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
	r.Get("/health", handler.Health)
	r.Post("/pdf", handler.PDF)

	addr := ":" + port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
