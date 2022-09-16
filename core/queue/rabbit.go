package queue

import (
	"fmt"
	"os"

	"bus-api/core/config"

	"github.com/streadway/amqp"
)

type Conn struct {
	Channel  *amqp.Channel
	Exchange string
}

// GetConn -
func GetConn() (Conn, error) {
	host := config.ReadConfig("queue.host")
	port := config.ReadConfig("queue.port")
	user := config.ReadConfig("queue.user")
	password := config.ReadConfig("queue.password")
	exchange := config.ReadConfig("queue.exchange")
	uri := "amqp://" + user + ":" + password + "@" + host + ":" + port + "/"
	conn, err := amqp.Dial(uri)
	if err != nil {
		return Conn{}, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return Conn{}, err
	}
	err = ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	if err != nil {
		return Conn{}, err
	}
	return Conn{
		Channel:  ch,
		Exchange: exchange,
	}, err
}

// Publish -
func (conn Conn) Publish(routingKey string, data []byte) error {
	return conn.Channel.Publish(
		conn.Exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         data,
			DeliveryMode: amqp.Persistent,
		})
}

// StartConsumer -
func (conn Conn) StartConsumer(queueName, routingKey string, handler func(d amqp.Delivery) bool) error {

	// create the queue if it doesn't already exist
	_, err := conn.Channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	// bind the queue to the routing key
	err = conn.Channel.QueueBind(queueName, routingKey, conn.Exchange, false, nil)
	if err != nil {
		return err
	}
	err = conn.Channel.Qos(4, 0, false)
	if err != nil {
		return err
	}

	msgs, err := conn.Channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			if handler(msg) {
				msg.Ack(false)
			} else {
				msg.Nack(false, true)
			}
		}
		fmt.Println("Rabbit consumer closed - critical Error")
		os.Exit(1)
	}()
	return nil
}
