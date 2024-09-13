package main

import (
	"context"
	"lengfengbyit/go-demos/rabbitmq/dle/helper"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// 1. 连接 mq
	conn := helper.ConnMQ()
	defer conn.Close()

	// 2. 获取通道
	ch, err := conn.Channel()
	helper.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	ch.Confirm(false)

	// 2.1 获取死信交换机
	dleExchange, _ := helper.CreateDealLetterQueue(ch)

	// 3. 声明队列
	ch.QueueDelete("world", false, false, false)
	queue, err := ch.QueueDeclare(
		"world",
		false,
		true,
		false,
		false,
		map[string]any{
			"x-dead-letter-exchange":    dleExchange,
			"x-dead-letter-routing-key": "",
			"x-message-ttl":             10000, // 消息过期时间，过期后就会自动到死信队列
		},
	)
	helper.FailOnError(err, "Failed to declare a queue")

	// 监听 Return 事件
	go func() {
		for n := range ch.NotifyReturn(make(chan amqp.Return, 1)) {
			log.Printf("return: %+v\n", n)
		}
	}()

	context.Context()

	// 发布消息
	err = ch.Publish(
		"",         // 交换机的名称，空字符串代表默认交换机
		queue.Name, // 路由键，默认交换机使用队列名做路由键
		true,       // 为 true : 当消息无法路由到任何队列时，RabbitMQ 会返回一个 return 事件给生产者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("hello world, this is dead letter"),
			UserId:      "myuser", // 必须和链接的登录名一致
			AppId:       "200",
			MessageId:   "300",
		},
	)
	helper.FailOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s\n", "hello world, this is dead letter")

	// 等待一段时间
	time.Sleep(10 * time.Second)
}
