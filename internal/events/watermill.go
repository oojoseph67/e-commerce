package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-aws/sqs"
	"github.com/ThreeDotsLabs/watermill/message"
	appConfig "github.com/oojoseph67/ecommerce/internal/config"
	"github.com/oojoseph67/ecommerce/internal/providers"
)

type EventPublisher struct {
	publisher message.Publisher
	queueName string
}

func NewEventPublisher(ctx context.Context, cfg *appConfig.AWSConfig) (*EventPublisher, error) {

	logger := watermill.NewStdLogger(true, true)

	// create aws config
	awsConfig, err := providers.CreateAwsConfig(ctx, cfg.S3Endpoint, cfg.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS config: %v", err)
	}

	// create watermill SQS publisher config
	publisherConfig := sqs.PublisherConfig{
		AWSConfig: awsConfig,
		Marshaler: nil,
	}

	// creating the publisher with custom config
	publisher, err := sqs.NewPublisher(publisherConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create publisher: %v", err)
	}

	event := &EventPublisher{
		publisher: publisher,
		queueName: cfg.EventQueueName,
	}

	return event, nil

}

func (ep *EventPublisher) Publish(eventType string, payload interface{}, metadata map[string]string) error {

	// converting our payload to json
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), data)

	// add metadata
	msg.Metadata.Set("event_type", eventType)
	for k, v := range metadata {
		msg.Metadata.Set(k, v)
	}

	return ep.publisher.Publish(ep.queueName, msg)
}

func (ep *EventPublisher) Close() error {
	return ep.publisher.Close()
}
