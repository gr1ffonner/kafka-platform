# Kafka vs NATS vs RabbitMQ Comparison

## Core Characteristics

| Feature | Kafka | NATS | RabbitMQ |
|---------|-------|------|----------|
| **Web UI** | Kafka UI - view topics, publish messages | Limited - basic monitoring only | Management UI (port 15672) - view queues/exchanges, publish messages |

у кафки клиента нет хелсчека
нужно разворачивать брокер для того чтобы начать принимать сообщения или продусить их
у рэббита похожая ситуация, нужно сначала создать exchange забиндить очереди

"amqp://guest:guest@localhost:5672/" - в кролик дснку нужно прокидывать протокол и креды
в кафку и натс по дефолту не нужно