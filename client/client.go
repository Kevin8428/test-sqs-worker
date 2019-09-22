package client

import "github.com/aws/aws-sdk-go/service/sqs"

type Client interface {
	ReceiveMessage() (*Message, error)
	SendMessage(payload interface{}) (string, error)
	DeleteMessage(msg *Message) error
	CreateQueue(inp *sqs.CreateQueueInput) (*sqs.CreateQueueOutput, error)
	ChangeMessageVisibilityTimeout(msg *Message, timeout int) error
	Ping() error
}
