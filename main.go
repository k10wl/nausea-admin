package main

import (
	"html/template"
	"log"
	"os"

	"nausea-admin/internal/cloudflare"
	"nausea-admin/internal/storage"
)

func main() {
	projectID := os.Getenv("PROJECT_ID")
	port := os.Getenv("PORT")

	db := NewDB(projectID)
	t := template.Must(template.ParseGlob("views/**"))

	c := cloudflare.NewClient()
	storage := storage.NewStorage(c)

	s := NewServer(":"+port, db, t, storage)
	if err := s.Run(); err != nil {
		log.Fatalf("FATAL SERVER ERROR: %v\n", err)
	}
}
