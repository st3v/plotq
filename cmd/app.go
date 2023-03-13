package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/st3v/plotq/filestore"
	"github.com/st3v/plotq/handler"
	"github.com/st3v/plotq/jobqueue"
	"github.com/st3v/plotq/spooler"
)

var (
	queueDir  = filepath.Join("data", "queue")
	uploadDir = filepath.Join("data", "upload")
)

func main() {
	queue, err := jobqueue.OpenLocal(queueDir)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create job queue: %w", err))
	}

	uploadStore, err := filestore.NewLocalStore(uploadDir)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create upload file store: %w", err))
	}

	handler := handler.New(spooler.NewSpooler(queue, uploadStore))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting service - http://localhost:%s/v1/docs\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), handler); err != nil {
		log.Fatal(err)
	}
}
