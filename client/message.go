package client

import "encoding/json"

type Metadata struct {
	ReceiveCount int
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
