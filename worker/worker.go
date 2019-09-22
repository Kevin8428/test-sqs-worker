package worker

import (
	"fmt"

	"github.com/kevin8428/test-sqs-worker/client"
)

type Worker struct {
	Handler                 Handler
	Queue                   client.Client
	BackoffExponent         int
	IgnoreDuplicateMessages bool
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
	err = w.Handler.Handle(message.Payload, message.MessageAttributes)
	if err != nil {
		fmt.Println("error handling message: ", err)
		return
	}
	err = w.Queue.DeleteMessage(message)
	if err != nil {
		fmt.Println("error deleting message: ", err)
		return
	}

}
