package hemmingway

import (
	"github.com/streadway/amqp"
	"log"
)


func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type RabbitWorker struct {
	Host string
	ExchangeName string
	ExchangeKind string
	QueueName string
	connection *amqp.Connection
	channel *amqp.Channel
	queue amqp.Queue
}

func (r *RabbitWorker) MakeConnection(){
	conn, err := amqp.Dial(r.Host)
	FailOnError(err, "Failed to connect to RabbitMQ")

	r.connection = conn
}

func (r *RabbitWorker) MakeChannel() {
	channel, err := r.connection.Channel()
	FailOnError(err, "Failed to open a channel")

	r.channel = channel
}

func (r *RabbitWorker) DeclareExchange() {
	err := r.channel.ExchangeDeclare(
		r.ExchangeName,
		r.ExchangeKind,
		false,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to open a channel")
}

func (r *RabbitWorker) DeclareQueue() {
	q, err := r.channel.QueueDeclare(
		r.QueueName, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	r.queue = q
}

func (r *RabbitWorker) BindQueueToExchange() {
	err := r.channel.QueueBind(
		r.QueueName, // queue name
		"",     // routing key
		r.ExchangeName, // exchange
		false,
		nil,
	)
	FailOnError(err, "Failed to bind queue")
}

func (r *RabbitWorker) SendPayload(payload []byte, contentType string, routingKey string) map[string]interface{}{
	err := r.channel.Publish(
		r.ExchangeName,     // exchange
		routingKey, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
			ContentType: contentType,
			Body: payload,
		})
	FailOnError(err, "Failed to publish a message")

	output := make(map[string]interface{})
	output["status"] = 200

	return output
}

func (r *RabbitWorker) Consume() {
	messages, err := r.channel.Consume(
		r.QueueName, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range messages {
			// TODO: Do something with messages here
			println(d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (r *RabbitWorker) CleanUp() {
	r.connection.Close()
	r.channel.Close()
}