package broker

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

var (
	topics  = []string{"test"}
	brokers = []string{"localhost:9093"}
)

// Broker handles Kafka infrastructure setup
type Broker struct {
	brokers []string
	topics  []string
	conn    *kafka.Conn
}

// NewBroker creates a new broker
func NewBroker(brokers []string) *Broker {
	b := &Broker{
		brokers: brokers,
		topics:  topics,
		conn:    nil,
	}

	// Initialize connection at struct level
	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		slog.Warn("Failed to establish initial connection", "error", err)
	} else {
		b.conn = conn
	}

	return b
}

// InitInfrastructure initializes all Kafka infrastructure
func (b *Broker) InitInfrastructure(ctx context.Context) error {
	slog.Info("Initializing Kafka infrastructure")

	// Ensure all topics exist
	for _, topicName := range b.topics {
		if err := b.ensureTopicExists(topicName); err != nil {
			return fmt.Errorf("failed to ensure topic '%s' exists: %w", topicName, err)
		}
	}

	slog.Info("Kafka infrastructure initialized successfully")
	return nil
}

// ensureTopicExists creates a topic if it doesn't exist
func (b *Broker) ensureTopicExists(topicName string) error {
	// Ensure we have a connection
	if b.conn == nil {
		conn, err := kafka.Dial("tcp", b.brokers[0])
		if err != nil {
			return fmt.Errorf("failed to connect to kafka: %w", err)
		}
		b.conn = conn
	}

	// Try to check if topic exists by reading partitions
	partitions, err := b.conn.ReadPartitions(topicName)
	if err != nil || len(partitions) == 0 {
		slog.Info("Topic does not exist, creating topic", "topic", topicName)

		// Create topic using existing connection
		topicConfig := kafka.TopicConfig{
			Topic:             topicName,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}

		err = b.conn.CreateTopics(topicConfig)
		if err != nil {
			return fmt.Errorf("failed to create topic '%s': %w", topicName, err)
		}

		slog.Info("Topic created successfully", "topic", topicName)
	} else {
		slog.Info("Topic already exists", "topic", topicName)
	}

	return nil
}

// Close closes the broker connection
func (b *Broker) Close() error {
	if b.conn != nil {
		return b.conn.Close()
	}
	return nil
}
