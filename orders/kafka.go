package main

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

const (
	KafkaServerAddress = "kafka-svc.default.svc.cluster.local:9092"
	KafkaTopic         = "orders"
)

func setupProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{KafkaServerAddress},
		config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup producer: %w", err)
	}
	return producer, nil
}

func sendKafkaMessage(producer sarama.SyncProducer, ctx *gin.Context, orderID, productID, quantity int64) error {
	notification := Notification{
		ProductID: productID,
		Quantity:  quantity,
	}
	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: KafkaTopic,
		Key:   sarama.StringEncoder(orderID),
		Value: sarama.StringEncoder(notificationJSON),
	}

	_, _, err = producer.SendMessage(msg)
	return err
}


