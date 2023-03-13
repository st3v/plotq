package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/st3v/plotq/converter"
	"github.com/st3v/plotq/filestore"
	"github.com/st3v/plotq/handler"
	"github.com/st3v/plotq/jobqueue"
	"github.com/st3v/plotq/spooler"
	"github.com/st3v/plotq/worker"
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

	converter := converter.Vpype()
	spool := spooler.NewSpooler(queue, uploadStore, converter.Convert)
	handler := handler.New(spool)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go worker.Run(ctx, spool)

	log.Printf("Starting service - http://localhost:%s/v1/docs\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), handler); err != nil {
		log.Fatal(err)
	}
}
