package sqsworker

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/kevin8428/test-sqs-worker/client"
	"github.com/kevin8428/test-sqs-worker/worker"
	"github.com/kevin8428/test-sqs-worker/workerpool"
)

type SQSConsumer struct {
	workers *workerpool.Pool
	worker  *worker.Worker
}

// SQS driver that satisfies the queue.Client interface
type SQSClientTwo struct {
	*sqs.SQS
	QueueURL string
}

type ConsumerConfig struct {
	// Handler                          worker.Handler
	// AWSRegion                        string
	// QueueURL                         string
	// RateLimit                        int
	// Parallel                         int
	// Pool                             int
	// SQSWorkerIgnoreDuplicateMessages bool
	// SQSWorkerRedisConnectionString   string
	// BackoffExponent                  int
	// BackoffStrategy                  string
}

func NewConsumer(handler worker.Handler) *SQSConsumer {
	awsSession := session.New()
	awsConfig := aws.NewConfig().WithRegion(os.Getenv("AWS_REGION"))
	queueClient := &client.SQSClient{
		QueueURL: os.Getenv("SQS_QUEUE_URL"),
		SQS:      sqs.New(awsSession, awsConfig),
	}
	pool := &workerpool.Pool{
		WaitTime: 100,
		Size:     10,
	}
	worker := &worker.Worker{
		Queue:   queueClient,
		Handler: handler,
	}
	return &SQSConsumer{
		worker:  worker,
		workers: pool,
	}
}

func (s *SQSConsumer) Start() {
	go s.workers.Start(s.worker)
}
