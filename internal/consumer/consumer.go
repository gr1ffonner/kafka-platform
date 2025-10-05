package consumer

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

var (
	topic         = "test"
	brokers       = []string{"localhost:9093"}
	consumerGroup = "test"
)

type Consumer struct {
	Reader        *kafka.Reader
	topic         string
	consumerGroup string
}

func NewConsumer(ctx context.Context) *Consumer {
	slog.Info("creating Kafka consumer", "topic", topic, "consumerGroup", consumerGroup)
	return &Consumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			// Brokers: list of Kafka broker addresses to connect to
			Brokers: brokers,

			// Topic: name of the Kafka topic to read messages from
			Topic: topic,

			// Consumer Group: enables offset management and acknowledgments
			// Messages are automatically acknowledged after processing
			GroupID: consumerGroup,

			// When using consumer groups, omit Partition field to read from all partitions
			// Kafka will automatically assign partitions to this consumer

			// MinBytes: minimum number of bytes to fetch per request
			// Larger values improve throughput but increase latency
			MinBytes: 10e3, // 10KB

			// MaxBytes: maximum number of bytes to fetch per request
			// Prevents memory issues with large messages
			MaxBytes: 10e6, // 10MB
		}),
		topic:         topic,
		consumerGroup: consumerGroup,
	}
}

// StartConsuming starts consuming messages from Kafka
func (c *Consumer) StartConsuming(ctx context.Context) error {
	slog.Info("starting Kafka consumer")

	for {
		select {
		case <-ctx.Done():
			slog.Info("consumer context cancelled, stopping")
			return ctx.Err()
		default:
			// Read message - this will block until a message is available or context is cancelled
			message, err := c.Reader.ReadMessage(ctx)

			if err != nil {
				if ctx.Err() != nil {
					slog.Info("consumer context cancelled, stopping")
					return ctx.Err()
				}
				slog.Error("failed to read message", "error", err)
				continue
			}

			// Process the message
			c.processMessage(message)
		}
	}
}

// processMessage processes a single Kafka message
func (c *Consumer) processMessage(message kafka.Message) {
	slog.Info("message received",
		"topic", message.Topic,
		"partition", message.Partition,
		"offset", message.Offset,
		"key", string(message.Key),
		"timestamp", message.Time,
		"size", len(message.Value),
	)

	// Try to parse as JSON
	var jsonMsg any
	if err := json.Unmarshal(message.Value, &jsonMsg); err == nil {
		slog.Info("message received", "message", jsonMsg)
	} else {
		slog.Warn("Message received is not valid JSON")
		slog.Info("Message received", "message", string(message.Value))
	}
}

func (c *Consumer) Close() error {
	return c.Reader.Close()
}
