package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/jjunary0425/aws-sqs-worker-001/pkg/internal/cloud"
)

type SqsClient struct {
	client *sqs.Client
	url    string
}

func NewSqSClient(url string) *SqsClient {
	config, err := New()
	if err != nil {
		fmt.Println("Error loading AWS config:", err)
		return nil
	}
	sqsClient := sqs.NewFromConfig(*config)
	return &SqsClient{
		client: sqsClient,
		url:    url,
	}
}

func (s *SqsClient) Send(ctx context.Context, req *cloud.SendRequest) (string, error) {
	attrs := make(map[string]types.MessageAttributeValue)

	for _, attr := range req.Attributes {
		attrs[attr.Key] = types.MessageAttributeValue{
			StringValue: aws.String(attr.Value),
			DataType:    aws.String(attr.Type),
		}
	}

	input := &sqs.SendMessageInput{
		QueueUrl:          aws.String(s.url),
		MessageBody:       aws.String("Hello, SQS!"),
		MessageAttributes: attrs,
	}

	res, err := s.client.SendMessage(context.TODO(), input)
	if err != nil {
		panic("failed to send message, " + err.Error())
	}

	return *res.MessageId, nil
}

func (s *SqsClient) Receive(ctx context.Context, queueURL string, maxMsg int64) ([]types.Message, error) {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(s.url),
		MaxNumberOfMessages: 10, // Number of messages to retrieve (up to 10)
		WaitTimeSeconds:     20, // Long polling timeout (in seconds)
	}

	resp, err := s.client.ReceiveMessage(context.TODO(), input)
	if err != nil {
		panic("failed to receive messages, " + err.Error())
	}

	for _, msg := range resp.Messages {
		fmt.Println("Message ID:", *msg.MessageId)
		fmt.Println("Message Body:", *msg.Body)
		// Process the message as needed
	}

	return resp.Messages, nil
}

func (s *SqsClient) Delete(ctx context.Context, queueURL, rcvHandle string) error {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.url),
		ReceiptHandle: aws.String("RECEIPT_HANDLE_OF_THE_MESSAGE_TO_DELETE"),
	}

	_, err := s.client.DeleteMessage(context.TODO(), input)
	if err != nil {
		panic("failed to delete message, " + err.Error())
	}
	return nil
}
