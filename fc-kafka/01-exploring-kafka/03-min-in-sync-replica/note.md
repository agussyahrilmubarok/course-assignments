# Min In Sync Replica

## Introduction

Apache Kafka is a **distributed event streaming platform** designed for high throughput, fault tolerance, and scalability.
It is commonly used for:

* **Event-driven architectures** (e.g., order events, user activity tracking).
* **Real-time data pipelines** (streaming from applications to databases, analytics systems, or other services).
* **Decoupling microservices** through asynchronous communication.

Kafka organizes data into **topics**, which are further divided into **partitions** and replicated across brokers. Replication ensures both parallelism and fault tolerance.

In this article, we will focus on the concept of **Minimum In-Sync Replicas (min.insync.replicas)** and why it is critical for durability and reliability in a Kafka cluster.

---

## What Is Min In-Sync Replicas?

A Kafka topic partition always has one **leader replica** and several **follower replicas**. Followers copy data from the leader to maintain redundancy. The set of replicas that are fully caught up with the leader is known as the **In-Sync Replica set (ISR)**.

The configuration `min.insync.replicas` specifies the **minimum number of replicas** (including the leader) that must acknowledge a write before it is considered successful.

If the number of available replicas falls below this threshold, Kafka will reject produce requests with an error (`NotEnoughReplicasException`).

---

## Why Is It Important?

1. **Durability**
   If `min.insync.replicas` is set too low (e.g., 1), a producer can succeed even if only the leader has written the message. If the leader crashes before followers catch up, that data can be lost.

2. **Trade-off Between Availability and Safety**

   * Lower values (e.g., 1) → higher availability, lower durability.
   * Higher values (e.g., 2 or more) → stronger durability guarantees, but may reduce availability if some brokers are offline.

3. **Recommended Setting**
   In a 3-broker cluster with `replication.factor=3`, it is common to set:

   * `min.insync.replicas=2`
     This ensures at least the leader plus one follower confirm the write, balancing safety and availability.

---

## Example Configurations

### Broker-Level Default

Add the following environment variable to each Kafka broker (e.g., in `docker-compose.yml`):

```yaml
environment:
  KAFKA_MIN_INSYNC_REPLICAS: 2
```

This sets the default for all newly created topics, unless overridden.

---

### Topic-Level Override

You can also set or update `min.insync.replicas` for a specific topic:

```bash
# Create topic with custom min.insync.replicas
kafka-topics --bootstrap-server kafka1:29092 \
  --create \
  --topic example.order.events \
  --partitions 3 \
  --replication-factor 3 \
  --config min.insync.replicas=2

# Alter existing topic
kafka-configs --bootstrap-server kafka1:29092 \
  --alter --entity-type topics --entity-name example.order.events \
  --add-config min.insync.replicas=2
```

---

## Producer Considerations

To take advantage of `min.insync.replicas`, producers should use the acknowledgment setting `acks=all`.

This ensures the producer only receives a success response if the leader and the required number of in-sync replicas have written the record.

If `acks=1` is used, the producer will only wait for the leader’s acknowledgment, bypassing the durability guarantees of `min.insync.replicas`.

---

## Conclusion

The `min.insync.replicas` setting plays a crucial role in defining Kafka’s durability guarantees:

1. It specifies how many replicas must acknowledge a write before it is considered successful.
2. It protects against data loss if a leader broker crashes.
3. It requires a balance between **availability** and **durability**, with `min.insync.replicas=2` being a common best practice for 3-replica clusters.

By combining `min.insync.replicas` with proper producer configurations (`acks=all`), you can build a more **resilient and fault-tolerant event-driven system**.

This forms the foundation for reliable event processing pipelines and microservice communication using Apache Kafka.
