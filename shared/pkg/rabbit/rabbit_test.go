package rabbit

import (
	"testing"
)

func TestNew_ReturnsNil_OnDialError(t *testing.T) {
	c := New("amqp://invalid-host:5672/")
	if c != nil {
		t.Fatalf("expected nil client for invalid url")
	}
}

func TestPublishJSON_NoOp_WhenClientNil(t *testing.T) {
	var c *Client
	if err := c.PublishJSON("events", "rk", map[string]any{"a": 1}); err != nil {
		t.Fatalf("expected nil error on nil client got %v", err)
	}
	// Also when channel is nil
	c2 := &Client{}
	if err := c2.PublishJSON("events", "rk", map[string]any{"a": 1}); err != nil {
		t.Fatalf("expected nil error on nil channel got %v", err)
	}
}

func TestConsume_NoOp_WhenClientNil(t *testing.T) {
	var c *Client
	msgs, err := c.Consume("events", "q", []string{"rk"})
	if err != nil || msgs != nil {
		t.Fatalf("expected nil,nil on nil client got msgs=%v err=%v", msgs, err)
	}
	// Also when channel is nil
	c2 := &Client{}
	msgs, err = c2.Consume("events", "q", []string{"rk"})
	if err != nil || msgs != nil {
		t.Fatalf("expected nil,nil on nil channel got msgs=%v err=%v", msgs, err)
	}
}

func TestClose_Safe_OnNil(t *testing.T) {
	var c *Client
	c.Close() // should not panic
	c2 := &Client{}
	c2.Close() // should not panic
}
