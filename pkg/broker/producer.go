package broker

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaWriter interface for dependency injection
type KafkaWriter interface {
	Init(ctx context.Context) error
	Close() error
	GetWriter() *kafka.Writer
}

type kafkaWriter struct {
	writer *kafka.Writer
}

func NewKafkaWriter() KafkaWriter {
	return &kafkaWriter{}
}

func (kw *kafkaWriter) Init(ctx context.Context) error {
	kw.writer = &kafka.Writer{
		// Addr: list of broker addresses to connect to
		Addr: kafka.TCP(brokers...),

		// Topic: default topic to write messages to (can be overridden per message)
		Topic: topics[0],

		// Balancer: algorithm for distributing messages across partitions
		// LeastBytes balances by choosing partition with least data
		Balancer: &kafka.LeastBytes{},

		// BatchTimeout: maximum time to wait before flushing a batch
		// Smaller values reduce latency but increase overhead
		BatchTimeout: 10 * time.Millisecond,

		// BatchSize: maximum number of messages to collect before flushing
		// Larger batches improve throughput but increase memory usage
		BatchSize: 100,
	}

	return nil
}

func (kw *kafkaWriter) Close() error {
	if kw.writer != nil {
		return kw.writer.Close()
	}
	return nil
}

func (kw *kafkaWriter) GetWriter() *kafka.Writer {
	return kw.writer
}
