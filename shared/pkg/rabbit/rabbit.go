package rabbit

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func New(url string) *Client {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil
	}
	return &Client{conn: conn, ch: ch}
}

func (c *Client) PublishJSON(exchange, routingKey string, payload map[string]any) error {
	if c == nil || c.ch == nil {
		return nil
	}
	if err := c.ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
		return err
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return c.ch.PublishWithContext(context.Background(), exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        b,
	})
}

func (c *Client) Consume(exchange, queue string, bindings []string) (<-chan amqp.Delivery, error) {
	if c == nil || c.ch == nil {
		return nil, nil
	}
	if err := c.ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
		return nil, err
	}
	q, err := c.ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	for _, key := range bindings {
		if e := c.ch.QueueBind(q.Name, key, exchange, false, nil); e != nil {
			return nil, e
		}
	}
	msgs, err := c.ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (c *Client) Close() {
	if c == nil {
		return
	}
	if c.ch != nil {
		_ = c.ch.Close()
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
