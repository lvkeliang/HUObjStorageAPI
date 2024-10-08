package rabbitmq

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	channel  *amqp.Channel
	Name     string
	exchange string
}

func New(s string) *RabbitMQ {
	conn, err := amqp.Dial(s)

	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	q, err := ch.QueueDeclare(
		"",
		false,
		true,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	mq := new(RabbitMQ)
	mq.channel = ch
	mq.Name = q.Name
	return mq
}

func (q *RabbitMQ) Bind(exchange string) {
	err := q.channel.QueueBind(
		q.Name,
		"", //路由键,留空表示可以接收交换机的所有消息
		exchange,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	q.exchange = exchange
}

func (q *RabbitMQ) Send(queue string, body interface{}) {
	str, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	err = q.channel.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body:    str,
		},
	)

	if err != nil {
		panic(err)
	}
}

func (q *RabbitMQ) Publish(exchange string, body interface{}) {
	str, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	err = q.channel.Publish(
		exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body:    str,
		},
	)

	if err != nil {
		panic(err)
	}
}

func (q *RabbitMQ) Consume() <-chan amqp.Delivery {
	c, err := q.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}
	return c
}

func (q *RabbitMQ) Close() {
	q.channel.Close()

}
