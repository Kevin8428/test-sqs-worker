package client

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/service/sqs"
)

type Metadata struct {
	ReceiveCount int
}

type Client interface {
	ReceiveMessage() (*Message, error)
	SendMessage(payload interface{}) (string, error)
	DeleteMessage(msg *Message) error
	CreateQueue(inp *sqs.CreateQueueInput) (*sqs.CreateQueueOutput, error)
	ChangeMessageVisibilityTimeout(msg *Message, timeout int) error
	Ping() error
}

type Message struct {
	Payload           []byte
	MessageID         string
	ReceiptHandle     string
	Metadata          Metadata
	MessageAttributes map[string]string
}

func (m *Message) UnmarshalInto(dest interface{}) error {
	return json.Unmarshal(m.Payload, dest)
}
