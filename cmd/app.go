package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/st3v/plotq/handler"
	"github.com/st3v/plotq/manager"
	"github.com/st3v/plotq/queue"
)

var queueDir = path.Join("data", "queue")

func main() {
	queue, err := queue.NewJobQueue(queueDir)
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
