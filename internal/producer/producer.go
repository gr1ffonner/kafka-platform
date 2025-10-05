package producer

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	topic   = "test"
	brokers = []string{"localhost:9093"}
)

type Producer struct {
	writer *kafka.Writer
	topic  string
}

func NewProducer(ctx context.Context) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			// Addr: list of broker addresses to connect to
			Addr: kafka.TCP(brokers...),

			// Topic: default topic to write messages to (can be overridden per message)
			Topic: topic,

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
		topic: topic,
	}
}

func (p *Producer) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}

// PublishMessage publishes a message synchronously
func (p *Producer) PublishMessage(ctx context.Context, message any) error {
	slog.Info("publishing message", "message", message, "topic", p.topic)

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
	}

	// Send message
	if err := p.writer.WriteMessages(ctx, kafkaMsg); err != nil {
		slog.Error("failed to publish message", "error", err)
		return err
	}

	slog.Info("message published", "message", message, "topic", p.topic)
	return nil
}
