package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/antonybholmes/go-sys/log"

	"github.com/antonybholmes/go-sys"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/panjf2000/ants"
)

type (
	SQSSender struct {
		Client   *sqs.Client
		QueueURL string
	}

	// Standardized consumer of sqs messages to act as foundation
	// for building specialized SQS consumers like email consumers
	SQSConsumer struct {
		Client       *sqs.Client
		pool         *ants.Pool
		backoff      *sys.Backoff
		QueueURL     string
		Threads      int
		WaitTimeSecs int
	}
)

// NewSQSSender creates a new instance of SQSSender with the given SQS client and queue URL.
func NewSQSSender(client *sqs.Client, queueURL string) *SQSSender {
	return &SQSSender{
		Client:   client,
		QueueURL: queueURL,
	}
}

func (s *SQSSender) Send(ctx context.Context, payload []byte) error {
	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(s.QueueURL),
		MessageBody: aws.String(string(payload))}

	if _, err := s.Client.SendMessage(ctx, input); err != nil {
		return fmt.Errorf("send SQS message: %w", err)
	}
	return nil
}

func (s *SQSSender) SendFifo(ctx context.Context,
	id string,
	groupID string,
	payload []byte) error {
	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(s.QueueURL),
		MessageBody: aws.String(string(payload))}

	if strings.HasSuffix(s.QueueURL, ".fifo") {
		input.MessageGroupId = aws.String(groupID)
		input.MessageDeduplicationId = aws.String(id)
	}

	if _, err := s.Client.SendMessage(ctx, input); err != nil {
		return fmt.Errorf("send SQS message: %w", err)
	}
	return nil
}

func NewSQSConsumer(client *sqs.Client,
	queueURL string,
	threads int,
	waitTimeSecs int,
) (*SQSConsumer, error) {

	// cfg, err := config.LoadDefaultConfig(context.Background())

	// if err != nil {
	// 	return nil, fmt.Errorf("unable to load AWS config: %w", err)
	// }

	// client := sqs.NewFromConfig(cfg)

	pool, err := ants.NewPool(threads)

	if err != nil {
		return nil, fmt.Errorf("unable to create worker pool with threads %d: %w", threads, err)
	}

	return &SQSConsumer{
		QueueURL:     queueURL,
		Threads:      threads,
		Client:       client,
		WaitTimeSecs: waitTimeSecs, // default wait time, can be adjusted
		pool:         pool,
		backoff:      sys.NewDefaultBackoff(),
	}, nil
}

// Consume starts consuming messages from the SQS queue.
// This is a long-running function that continuously polls the queue for messages.
func (c *SQSConsumer) Consume(ctx context.Context,
	processMessage func(ctx context.Context, message sqstypes.Message) error) {

	log.Debug().Msgf("start sqs %s", c.QueueURL)

	for {
		// pull as many messages as we have worker threads
		messageBatch, err := c.Client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(c.QueueURL),
			MaxNumberOfMessages: int32(c.Threads),
			WaitTimeSeconds:     int32(c.WaitTimeSecs),
		})

		if err != nil {
			log.Error().Msgf("failed to receive messages: %v", err)
			c.backoff.Sleep()
			continue
		}

		// process each message individually using the worker pool
		for _, message := range messageBatch.Messages {
			submitErr := c.pool.Submit(func() {

				messageID := GetMessageID(message)

				if message.Body == nil || message.ReceiptHandle == nil {
					log.Error().Msgf("message %s has no body or receipt handle", messageID)
					return
				}

				// Seems like have a message so pass off to the message processor
				err := processMessage(ctx, message)

				if err != nil {
					log.Error().Msgf("error processing SQS message: %v", err)
					return
				}

				// delete the message from the queue if processed successfully
				if err := c.DeleteMessage(ctx, message); err != nil {
					log.Error().Msgf("error deleting SQS message: %v", err)
				}
			})

			if submitErr != nil {
				log.Error().Msgf("failed to submit SQS message: %v", submitErr)
			}
		}
	}
}

func (c *SQSConsumer) DeleteMessage(ctx context.Context, message sqstypes.Message) error {
	if _, err := c.Client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.QueueURL),
		ReceiptHandle: message.ReceiptHandle,
	}); err != nil {
		return fmt.Errorf("delete message %s: %w", *message.MessageId, err)
	}

	return nil
}

func (c *SQSConsumer) Close() {
	c.pool.Release()
}

// GetMessageID returns the ID of the SQS message, or "<unknown>" if it is not set.
func GetMessageID(message sqstypes.Message) string {
	if message.MessageId != nil {
		return *message.MessageId
	}
	return "<unknown>"
}
