package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/st3v/plotq/handler"
	"github.com/st3v/plotq/jobqueue"
	"github.com/st3v/plotq/manager"
)

var (
	queueDir = filepath.Join("data", "queue")
)

func main() {
	queue, err := jobqueue.NewLocalQueue(queueDir)
	if err != nil {
		log.Fatal(err)
	}

	handler := handler.New(manager.NewJobManager(queue))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting service - http://localhost:%s/v1/docs\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), handler); err != nil {
		log.Fatal(err)
	}
}
