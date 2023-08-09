package tools

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

var kafkaConn *kafka.Writer

// to produce messages
func PushMsgToKafka(msg string) {
	err := kafkaConn.WriteMessages(
		context.Background(),
		kafka.Message{Value: []byte(msg)},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
}

func init() {
	kafkaConn = kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "my-topic",
	})
}
