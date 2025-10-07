package producer

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	brokers = []string{"localhost:9093"}
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(ctx context.Context) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			// Addr: list of broker addresses to connect to
			Addr: kafka.TCP(brokers...),

			// Balancer: algorithm for distributing messages across partitions
			// LeastBytes balances by choosing partition with least data
			Balancer: &kafka.LeastBytes{},

			// BatchTimeout: maximum time to wait before flushing a batch
			// Smaller values reduce latency but increase overhead
			BatchTimeout: 10 * time.Millisecond,

			// BatchSize: maximum number of messages to collect before flushing
			// Larger batches improve throughput but increase memory usage
			BatchSize: 100,
		},
	}
}

func (p *Producer) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}

// PublishMessage publishes a message synchronously
func (p *Producer) PublishMessage(ctx context.Context, topic string, message any) error {
	if topic == "" {
		return errors.New("topic is required")
	}
	slog.Info("publishing message", "message", message, "topic", topic)

	var messageBytes []byte

	// Serialize message to JSON
	messageBytes, err := json.Marshal(message)
	if err != nil {
		slog.Error("failed to marshal message", "error", err)
		return err
	}

	// Create Kafka message
	kafkaMsg := kafka.Message{
		Value: messageBytes,
		Topic: topic,
	}

	// Send message
	if err := p.writer.WriteMessages(ctx, kafkaMsg); err != nil {
		slog.Error("failed to publish message", "error", err)
		return err
	}

	slog.Info("message published", "message", message, "topic", topic)
	return nil
}
