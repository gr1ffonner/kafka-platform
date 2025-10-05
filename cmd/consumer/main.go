package main

import (
	"context"
	"kafka-platform/internal/consumer"
	"kafka-platform/pkg/broker"
	"kafka-platform/pkg/config"
	"kafka-platform/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger.InitLogger(cfg.Logger)
	logger := slog.Default()

	logger.Info("config and logger initialized")

	// Initialize Kafka broker
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize Kafka infrastructure
	kafkaBroker := broker.NewBroker([]string{cfg.KafkaDSN})
	defer kafkaBroker.Close()

	// Initialize infrastructure (create topics, etc.)
	if err := kafkaBroker.InitInfrastructure(ctx); err != nil {
		logger.Error("failed to initialize Kafka infrastructure", "error", err)
		os.Exit(1)
	}

	// Initialize consumer
	kafkaConsumer := consumer.NewConsumer(ctx)
	if kafkaConsumer == nil {
		logger.Error("failed to initialize Kafka consumer", "error", err)
		os.Exit(1)
	}
	defer kafkaConsumer.Close()

	// Setup signal handling for graceful shutdown
	signalCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	consumerDone := make(chan error, 1)
	go func() {
		consumerDone <- kafkaConsumer.StartConsuming(signalCtx)
	}()

	// Wait for shutdown signal or consumer error
	select {
	case <-signalCtx.Done():
		logger.Info("shutting down consumer")
		cancel() // Cancel the main context
	case err := <-consumerDone:
		if err != context.Canceled {
			logger.Error("consumer stopped with error", "error", err)
		}
	}

	logger.Info("consumer shutdown complete")
}
