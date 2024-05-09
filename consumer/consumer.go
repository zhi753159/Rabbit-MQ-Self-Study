package main

import (
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

	msgs, err := ch.Consume(
		q.Name,
		"c001",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		failOnError(err, "Failed to register a consumer")
		return
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			dt := time.Now()
			log.Printf("[%s] Received a messge: %s", dt.Format("2006-01-02 15:04:05"), d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
