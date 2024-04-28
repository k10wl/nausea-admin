package main

import (
	"log"
	"os"

	"nausea-admin/internal/cloudflare"
	"nausea-admin/internal/db"
	"nausea-admin/internal/firestore"
	"nausea-admin/internal/server"
	"nausea-admin/internal/storage"
)

func main() {
	projectID := os.Getenv("PROJECT_ID")
	port := os.Getenv("PORT")

	f := firestore.NewFirestoreClient(projectID)
	db := db.NewDB(f)

	c := cloudflare.NewClient()
	storage := storage.NewStorage(c)

	s := server.NewServer(":"+port, db, storage)
	if err := s.Run(); err != nil {
		log.Fatalf("FATAL SERVER ERROR: %v\n", err)
	}
}
