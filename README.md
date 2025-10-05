# Kafka Messaging Project

A Go project demonstrating Kafka messaging patterns including pub/sub, load balancing, and persistent messaging.

## Project Overview

This project showcases Kafka messaging capabilities with:
- **Publishers**: Message producers with various delivery guarantees
- **Consumers**: Message consumers with group coordination
- **Load Balancing**: Consumer groups for distributed processing
- **Persistence**: Kafka's built-in message persistence and replication
- **Monitoring**: Kafka cluster monitoring and metrics

## Services & Ports

| Service | Port | Description |
|---------|------|-------------|
| **Kafka Broker** | 9092 | Message broker |
| **Zookeeper** | 2181 | Kafka coordination service |
| **Schema Registry** | 8081 | Schema management |

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+
- Make

### Commands

#### Infrastructure Management
```bash
# Start Kafka cluster
make up

# Stop Kafka cluster
make down
```

#### Application
```bash
# Run Kafka consumer
make run-consumer

# Run Kafka producer
make run-producer
```

## Monitoring
- **Kafka UI**: http://localhost:8080
- **Schema Registry UI**: http://localhost:8081