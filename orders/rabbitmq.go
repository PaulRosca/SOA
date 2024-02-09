package main

import (
	"fmt"

	"github.com/streadway/amqp"
)

func sendRabbitMessage(channel *amqp.Channel, message string) error {
	// publishing a message
	err := channel.Publish(
		"",        // exchange
		"emails", // key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	return err
}

func connectToRabbit() (*amqp.Connection, *amqp.Channel) {
	connection, err := amqp.Dial("amqp://guest:guest@rabbitmq.default.svc.cluster.local:5672/")
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to RabbitMQ instance")

	// opening a channel over the connection established to interact with RabbitMQ
	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}

	// declaring queue with its properties over the the channel opened
	queue, err := channel.QueueDeclare(
		"emails", // name
		false,     // durable
		false,     // auto delete
		false,     // exclusive
		false,     // no wait
		nil,       // args
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Queue status:", queue)
	fmt.Println("Successfully published message")

	return connection, channel
}
