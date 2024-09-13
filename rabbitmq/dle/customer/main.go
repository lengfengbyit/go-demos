package main

import (
	"lengfengbyit/go-demos/rabbitmq/dle/helper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// 消费死信队列
func main() {

	// 1. 获取连接
	conn := helper.ConnMQ()
	defer conn.Close()

	// 2. 获取 channel
	ch, err := conn.Channel()
	helper.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 3. 获取死信队列名
	_, dleQueue := helper.CreateDealLetterQueue(ch)

	consume, err := ch.Consume(
		dleQueue,
		"",    // 消费者标签
		false, // 是否自动确认
		false, // 是否独占
		false, // 是否不是本地
		false, // 是否不等待
		nil,   //  其他参数
	)
	helper.FailOnError(err, "Failed to register a consumer")

	go func() {
		for msg := range consume {
			log.Println("[DLE] Received a message:", string(msg.Body), map[string]string{
				"UserID":    msg.UserId,
				"AppID":     msg.AppId,
				"MessageID": msg.MessageId,
			})
			_ = msg.Ack(false)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
}
