package client

import (
	"encoding/json"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQS driver that satisfies the queue.Client interface
type SQSClient struct {
	*sqs.SQS
	QueueURL string
}

func (s *SQSClient) SendMessage(payload interface{}) (string, error) {
	marshaled, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	input := &sqs.SendMessageInput{
		MessageBody: aws.String(string(marshaled)),
		QueueUrl:    aws.String(s.QueueURL),
	}

	output, err := s.SQS.SendMessage(input)
	if err != nil {
		return "", err
	}

	return *output.MessageId, nil
}

func (s *SQSClient) DeleteMessage(msg *Message) error {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.QueueURL),
		ReceiptHandle: aws.String(msg.ReceiptHandle),
	}

	_, err := s.SQS.DeleteMessage(input)

	return err
}

func (s *SQSClient) ReceiveMessage() (*Message, error) {
	input := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(s.QueueURL),

		// This defaults to 1 but let's be explicit
		MaxNumberOfMessages: aws.Int64(1),

		AttributeNames:        []*string{aws.String("All")},
		MessageAttributeNames: []*string{aws.String("All")},
	}

	output, err := s.SQS.ReceiveMessage(input)
	if err != nil {
		return nil, err
	}

	if len(output.Messages) > 0 {
		sqsMsg := output.Messages[0]

		receiveCount, ok := sqsMsg.Attributes["ApproximateReceiveCount"]
		receiveCountInt := int64(0)
		if ok {
			receiveCountInt, err = strconv.ParseInt(*receiveCount, 10, 64)
			if err != nil {
				return nil, err
			}
		}

		messageAttributes := map[string]string{}
		for key, value := range sqsMsg.MessageAttributes {
			if value.StringValue != nil {
				messageAttributes[key] = *value.StringValue
			}
		}

		msg := &Message{
			Payload:           []byte(*sqsMsg.Body),
			MessageID:         *sqsMsg.MessageId,
			ReceiptHandle:     *sqsMsg.ReceiptHandle,
			MessageAttributes: messageAttributes,
		}

		msg.Metadata.ReceiveCount = int(receiveCountInt)

		return msg, nil
	}

	return nil, nil
}

func (s *SQSClient) ChangeMessageVisibilityTimeout(msg *Message, timeout int) error {
	_, err := s.SQS.ChangeMessageVisibility(&sqs.ChangeMessageVisibilityInput{
		QueueUrl:          aws.String(s.QueueURL),
		ReceiptHandle:     aws.String(msg.ReceiptHandle),
		VisibilityTimeout: aws.Int64(int64(timeout)),
	})
	return err
}

// This just checks to make sure the config is all set up correctly. This could really be
// any function that makes a request to SQS.
func (s *SQSClient) Ping() error {
	req := new(sqs.GetQueueAttributesInput)
	req.QueueUrl = aws.String(s.QueueURL)

	_, err := s.SQS.GetQueueAttributes(req)
	return err
}
