package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func connectRabbitMQ() (*amqp.Connection, error) {
	username := os.Getenv("RABBIT_MQ_USERNAME")
	password := os.Getenv("RABBIT_MQ_PASSWORD")
	host := os.Getenv("RABBIT_MQ_HOST")
	port := os.Getenv("RABBIT_MQ_PORT")
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port)
	fmt.Print(connStr)
	conn, err := amqp.Dial(connStr)

	return conn, err
}

func main() {
	conn, err := connectRabbitMQ()
	if err != nil {
		failOnError(err, "Failed to connect to RabbitMQ")
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		failOnError(err, "Fail to open a channel")
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"chat",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World"
	err = ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text.plain",
			Body:        []byte(body),
		})
	if err != nil {
		failOnError(err, "Fail to publish message")
		return
	}
	dt := time.Now()
	log.Printf("[%s] Sent %s\n", dt.Format("2006-01-02 15:04:05"), body)
}
