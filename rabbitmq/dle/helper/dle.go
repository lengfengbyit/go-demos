package helper

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// dead-letter-queue 死信队列 demo

func ConnMQ() *amqp.Connection {
	conn, err := amqp.Dial("amqp://myuser:mypass@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	return conn
}

func CreateDealLetterQueue(ch *amqp.Channel) (string, string) {
	dleExchange := "dle_exchange"
	dleQueue := "dle_queue"

	// 创建死信交换机
	err := ch.ExchangeDeclare(
		dleExchange,
		amqp.ExchangeDirect,
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to DLE declare exchange")

	// 创建死信队列
	_, err = ch.QueueDeclare(
		dleQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to DLE declare queue")

	// 绑定死信队列到交换机
	err = ch.QueueBind(
		dleQueue,
		"",
		dleExchange,
		false,
		nil,
	)
	FailOnError(err, "Failed to DLE bind queue")

	log.Printf("[INFO] DLE queue created")

	return dleExchange, dleQueue
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s:%s", msg, err)
	}
}
