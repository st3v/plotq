package main

import (
	"log"
	"net/http"
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

	log.Println("Starting service - http://localhost:8080/v1/docs")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
