package worker

import (
	"fmt"

	"github.com/kevin8428/test-sqs-worker/client"
)

type Worker struct {
	Handler Handler
	// Metrics                 metrics.Client
	// BackoffStrategy         backoffStrategy
	Queue                   client.Client
	BackoffExponent         int
	IgnoreDuplicateMessages bool
	// RedisClient             *redis.Clientj
}

func (w *Worker) Work() {
	fmt.Println("requesting message from SQS")
	message, err := w.Queue.ReceiveMessage()
	if err != nil {
		fmt.Println("AWS SQS SDK error: ", err)
		return
	}
	if message == nil {
		fmt.Println("no message received")
		return
	}
	// to be run by service importing
	// err = w.Handler.Handle(message.Payload, message.Attributes)
}
