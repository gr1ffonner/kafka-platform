package main

import (
	"context"
	"kafka-platform/internal/producer"
	"kafka-platform/pkg/broker"
	"kafka-platform/pkg/config"
	"kafka-platform/pkg/logger"
	"log/slog"
	"os"
	"time"
)

const (
	message = "Hello, world"
	topic   = "test"
)

type Message struct {
	Message string `json:"message"`
}

func main() {
	msg := Message{
		Message: message,
	}

	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		panic(err)
	}
	slog.Info("Config loaded", "config", cfg)

	logger.InitLogger(cfg.Logger)
	logger := slog.Default()

	logger.Info("Config and logger initialized")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initialize Kafka infrastructure
	kafkaBroker := broker.NewBroker([]string{cfg.KafkaDSN})
	defer kafkaBroker.Close()

	// Initialize infrastructure (create topics, etc.)
	if err := kafkaBroker.InitInfrastructure(ctx); err != nil {
		logger.Error("Failed to initialize Kafka infrastructure", "error", err)
		os.Exit(1)
	}

	kafkaProducer := producer.NewProducer(ctx)
	defer kafkaProducer.Close()

	err = kafkaProducer.PublishMessage(ctx, topic, message)
	if err != nil {
		logger.Error("Failed to publish message", "error", err, "message", message)
		os.Exit(1)
	}

	err = kafkaProducer.PublishMessage(ctx, topic, msg)
	if err != nil {
		logger.Error("Failed to publish message", "error", err, "message", msg)
		os.Exit(1)
	}

	logger.Info("Successfully published all messages")
}
