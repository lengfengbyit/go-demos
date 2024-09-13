package main

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://myuser:mypass@localhost:5672/?connection_timeout=600000")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	ch.Confirm(false)
	go func() {
		// 确认监听，确认消息是否发送成功
		for confirmed := range ch.NotifyPublish(make(chan amqp.Confirmation, 1)) {
			log.Printf("Confirmation of delivery with delivery tag %v", confirmed)
		}
	}()

	go func() {
		// 监听消息是否投递到队列中
		for n := range ch.NotifyReturn(make(chan amqp.Return, 1)) {
			log.Printf("Confirmation of delivery with delivery tag %s", string(n.Body))
		}
	}()

	println(q.Name)

	msg := "Hello World"
	err = ch.Publish(
		"", // exchange
		//q.Name, // routing key
		"test",
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(msg),
			DeliveryMode: 2,
			UserId:       "myuser1",
		})

	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", msg)

	time.Sleep(10 * time.Second)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s:%s", msg, err)
	}
}
