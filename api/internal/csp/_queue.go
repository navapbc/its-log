package csp

import (
	"log"

	"github.com/navapbc/its-log/internal/base"
)

func Queue(op chan base.Job) {
	var counter int64 = 0
	for {
		next := <-op
		log.Printf("Queue[%d] %s\n", counter, next.Id)
		next.Op()
		next.Resp <- counter
		counter += 1
	}
}
