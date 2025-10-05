# Kafka Platform

## Services & Ports

| Service | Port | Description |
|---------|------|-------------|
| **Kafka Broker** | 9093 | Message broker |
| **Kafka UI** | 8881 | Kafka Web UI |
| **Zookeeper** | 2181 | Kafka coordination service |

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
- **Kafka UI**: http://localhost:8881