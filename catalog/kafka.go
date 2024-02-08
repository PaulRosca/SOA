package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	kafka "github.com/segmentio/kafka-go"
)

const (
	ConsumerGroup      = "orders-group"
	ConsumerTopic      = "orders"
	KafkaServerAddress = "kafka-svc.default.svc.cluster.local:9092"
)

type Notification struct {
	ProductID int64 `json:"productID"`
	Quantity  int64 `json:"quantity"`
}

func setupConsumer() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{KafkaServerAddress},
		GroupID:  ConsumerGroup,
		Topic:    ConsumerTopic,
		MaxBytes: 10e6, // 10MB
	})
	defer r.Close()
	ctx := context.Background()
	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			break
		}

		notification := Notification{}
		if err := json.Unmarshal(m.Value, &notification); err != nil {
			fmt.Println("Error", err.Error())
		}

		_, err = DB.Exec("UPDATE product SET stock = stock - ? WHERE id = ?", notification.Quantity, notification.ProductID)
		if err != nil {
			fmt.Println("Error", err.Error())
		}

		if err := r.CommitMessages(ctx, m); err != nil {
			log.Fatal("failed to commit messages:", err)
		}
	}
}
