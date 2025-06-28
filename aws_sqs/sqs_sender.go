package aws_sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"log"
	"os"

	"github.com/Tabernol/krasiot-sensor/model"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SqsNotifier struct {
	client   *sqs.Client
	queueURL string
}

func NewSqsNotifier(ctx context.Context, queueURL string) *SqsNotifier {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		log.Fatalf("❌ Failed to load AWS config: %v", err)
	}

	client := sqs.NewFromConfig(cfg)
	return &SqsNotifier{
		client:   client,
		queueURL: queueURL,
	}
}

func (n *SqsNotifier) SendEnrichedSensorData(ctx context.Context, data model.EnrichedSensorData) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	_, err = n.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &n.queueURL,
		MessageBody: aws.String(string(payload)),
	})
	return err
}
