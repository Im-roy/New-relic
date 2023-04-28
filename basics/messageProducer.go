package basics

import (
	"context"
	"log"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/segmentio/kafka-go"
)

func GetKafkaConnection() *kafka.Conn {
	topic := "my-topic"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10*time.Second))
	return conn
}

func WriteMessages(conn *kafka.Conn) {
	// to produce messages
	_, err := conn.WriteMessages(
		kafka.Message{Value: []byte("one!")},
		kafka.Message{Value: []byte("two!")},
		kafka.Message{Value: []byte("three!")},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

func MessageProducerSegment(ctx context.Context) newrelic.MessageProducerSegment {
	txn := newrelic.FromContext(ctx)
	s := newrelic.MessageProducerSegment{
		StartTime:            txn.StartSegmentNow(),
		Library:              "Kafka",
		DestinationType:      newrelic.MessageTopic,
		DestinationName:      "my-topic",
		DestinationTemporary: true,
	}
	return s
}

func InstrumentMessageProducer(ctx context.Context) {
	seg := MessageProducerSegment(ctx)
	defer seg.End()

	conn := GetKafkaConnection()
	WriteMessages(conn)
}