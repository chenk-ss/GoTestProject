package tools

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func init() {
	client := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"},
		GroupID:        "my-group",
		Topic:          "my-topic",
		Partition:      0,
		MaxBytes:       10e6,        // 10MB
		CommitInterval: time.Second, // flushes commits to Kafka every second
	})

	go func() {
		for {
			m, err := client.ReadMessage(context.Background())
			if err != nil {
				break
			}
			fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
			if err := client.CommitMessages(context.Background(), m); err != nil {
				log.Fatal("failed to commit messages:", err)
			}
		}

		if err := client.Close(); err != nil {
			log.Fatal("failed to close reader:", err)
		}
	}()
}
