# Kafka Load Balancing Documentation

## Overview

Kafka provides sophisticated load balancing through partitioning and consumer groups. This document covers how Kafka distributes messages across multiple consumers for efficient processing and horizontal scaling.

## Partitioning and Load Balancing

### How Kafka Distributes Work

```mermaid
graph TB
    A[Producer] -->|Publish Messages| B[Topic: orders]
    B --> C[Partition 0]
    B --> D[Partition 1]
    B --> E[Partition 2]
    B --> F[Partition 3]
    
    G[Consumer Group: processors] --> C
    G --> D
    G --> E
    G --> F
    
    C --> H[Consumer 1]
    D --> I[Consumer 2]
    E --> J[Consumer 3]
    F --> K[Consumer 4]
```

**Key Principles:**
- Each partition can only be consumed by one consumer in a group
- Messages within a partition are processed in order
- Load is distributed evenly across available consumers
- Adding more consumers than partitions provides no benefit

## Consumer Group Coordination

### Consumer Group Rebalancing

```mermaid
sequenceDiagram
    participant C1 as Consumer 1
    participant C2 as Consumer 2
    participant C3 as Consumer 3
    participant CG as Consumer Group
    participant K as Kafka Broker
    
    C1->>K: Join Group
    C2->>K: Join Group
    C3->>K: Join Group
    
    K->>CG: Assign Partitions
    Note over CG: Partition 0 -> Consumer 1<br/>Partition 1 -> Consumer 2<br/>Partition 2 -> Consumer 3
    
    C2->>K: Leave Group
    K->>CG: Rebalance Required
    Note over CG: Partition 1 -> Consumer 1 or 3
    
    C1->>K: Reassign Partitions
    C3->>K: Reassign Partitions
```

### Rebalancing Triggers

1. **Consumer Joins**: New consumer joins the group
2. **Consumer Leaves**: Consumer leaves or fails
3. **Topic Changes**: New partitions added to topic
4. **Group Membership Changes**: Any membership modification

## Load Balancing Strategies

### 1. Partition-Based Distribution

```mermaid
graph LR
    A[Topic: orders<br/>4 Partitions] --> B[Consumer Group]
    
    B --> C[Consumer 1<br/>Partition 0]
    B --> D[Consumer 2<br/>Partition 1]
    B --> E[Consumer 3<br/>Partition 2]
    B --> F[Consumer 4<br/>Partition 3]
    
    style C fill:#e1f5fe
    style D fill:#f3e5f5
    style E fill:#e8f5e8
    style F fill:#fff3e0
```

**Benefits:**
- Perfect load distribution (1 partition per consumer)
- Maximum parallelism
- Predictable processing patterns

### 2. Uneven Distribution Handling

```mermaid
graph LR
    A[Topic: orders<br/>3 Partitions] --> B[Consumer Group]
    
    B --> C[Consumer 1<br/>Partition 0]
    B --> D[Consumer 2<br/>Partition 1]
    B --> E[Consumer 3<br/>Partition 2]
    B --> F[Consumer 4<br/>Idle]
    
    style F fill:#ffebee
```

**Scenario:** More consumers than partitions
- Consumer 4 remains idle
- No additional processing benefit
- Consider scaling partitions instead

### 3. Over-Subscribed Partitions

```mermaid
graph LR
    A[Topic: orders<br/>2 Partitions] --> B[Consumer Group]
    
    B --> C[Consumer 1<br/>Partition 0]
    B --> D[Consumer 2<br/>Partition 1]
    B --> E[Consumer 3<br/>Partition 0<br/>Shared Load]
    
    style C fill:#e1f5fe
    style D fill:#f3e5f5
    style E fill:#ffebee
```

**Note:** This scenario is not possible in Kafka - each partition is consumed by exactly one consumer per group.

## Platform Implementation

### Consumer Configuration

```go
type Consumer struct {
    reader                *kafka.Reader
    logger                *slog.Logger
    handleMessagesCommand message.Command[message.HandleMessagesCommandOptions]
    config                *config.KafkaConsumer
    consumerID            string
}
```

**Load Balancing Features:**
- **Consumer ID**: Unique identifier for each consumer instance
- **Group ID**: Shared identifier for consumer group coordination
- **Partition Assignment**: Automatic partition assignment by Kafka
- **Batch Processing**: Configurable batch sizes for efficiency

### Batch Processing Load Balancing

```mermaid
graph TD
    A[Consumer 1<br/>Partition 0] --> B[Batch Size: 100]
    C[Consumer 2<br/>Partition 1] --> D[Batch Size: 100]
    E[Consumer 3<br/>Partition 2] --> F[Batch Size: 100]
    
    B --> G[Process 100 Messages]
    D --> H[Process 100 Messages]
    F --> I[Process 100 Messages]
    
    G --> J[Commit Offsets]
    H --> K[Commit Offsets]
    I --> L[Commit Offsets]
```

**Configuration:**
```json
{
  "readBatchSize": 100,
  "flushInterval": 1000,
  "flushMaxRetries": 3
}
```

## Scaling Strategies

### Horizontal Scaling

```mermaid
graph TB
    A[Initial: 1 Consumer] --> B[Scale: 2 Consumers]
    B --> C[Scale: 4 Consumers]
    
    A1[Consumer 1<br/>All Partitions] --> A
    
    B1[Consumer 1<br/>Partition 0,1] --> B
    B2[Consumer 2<br/>Partition 2,3] --> B
    
    C1[Consumer 1<br/>Partition 0] --> C
    C2[Consumer 2<br/>Partition 1] --> C
    C3[Consumer 3<br/>Partition 2] --> C
    C4[Consumer 4<br/>Partition 3] --> C
```

**Scaling Rules:**
1. **Optimal**: Number of consumers = Number of partitions
2. **Maximum**: Can have more consumers than partitions (some idle)
3. **Minimum**: Must have at least 1 consumer per partition for processing

### Partition Scaling

```mermaid
graph LR
    A[Topic: orders<br/>2 Partitions] --> B[Scale to 4 Partitions]
    B --> C[Scale to 8 Partitions]
    
    A1[Consumer 1: P0] --> A
    A2[Consumer 2: P1] --> A
    
    B1[Consumer 1: P0] --> B
    B2[Consumer 2: P1] --> B
    B3[Consumer 3: P2] --> B
    B4[Consumer 4: P3] --> B
```

**Partition Scaling Benefits:**
- Enables more parallel consumers
- Better load distribution
- Higher throughput capacity

## Load Balancing Scenarios

### Scenario 1: Perfect Distribution

```
Topic: orders (4 partitions)
Consumer Group: order-processors (4 consumers)

Consumer 1 -> Partition 0
Consumer 2 -> Partition 1
Consumer 3 -> Partition 2
Consumer 4 -> Partition 3

Result: Optimal load balancing
```

### Scenario 2: Consumer Failure

```
Initial:
Consumer 1 -> Partition 0, 1
Consumer 2 -> Partition 2, 3

Consumer 1 Fails:
Consumer 2 -> Partition 0, 1, 2, 3

Result: Automatic failover, temporary load increase
```

### Scenario 3: Consumer Recovery

```
During Failure:
Consumer 2 -> Partition 0, 1, 2, 3

Consumer 1 Recovers:
Consumer 1 -> Partition 0, 1
Consumer 2 -> Partition 2, 3

Result: Automatic rebalancing, load redistribution
```

### Scenario 4: New Consumer Joins

```
Initial:
Consumer 1 -> Partition 0, 1, 2, 3

New Consumer 2 Joins:
Consumer 1 -> Partition 0, 1
Consumer 2 -> Partition 2, 3

Result: Load automatically redistributed
```

## Performance Considerations

### Consumer Lag Monitoring

```mermaid
graph LR
    A[Producer<br/>High Throughput] --> B[Topic Partitions]
    B --> C[Consumer Group<br/>Slow Processing]
    
    D[Consumer Lag<br/>Messages Behind] --> C
    E[Scale Consumers<br/>Add More Instances] --> C
    F[Scale Partitions<br/>Increase Parallelism] --> B
```

**Lag Indicators:**
- High consumer lag = processing bottleneck
- Solution: Add more consumers or partitions
- Monitor through Kafka UI or metrics

### Batch Size Optimization

```go
// Small batches: Lower latency, higher overhead
config.ReadBatchSize = 10

// Large batches: Higher throughput, higher latency
config.ReadBatchSize = 1000

// Optimal: Balance based on message size and processing time
config.ReadBatchSize = 100
```

## Monitoring Load Balancing

### Kafka UI Monitoring

Access at:
- **Producer**: http://localhost:8888
- **Consumer**: http://localhost:8881

**Key Metrics:**
- Consumer group status
- Partition assignments
- Consumer lag per partition
- Message throughput rates

### Health Checks

```go
// Consumer health check
func (c *Consumer) HealthCheck() error {
    // Check Kafka connectivity
    // Verify consumer group membership
    // Monitor processing lag
}
```

**Health Indicators:**
- Consumer group membership
- Partition assignment status
- Processing lag thresholds
- Error rates and retry counts