package apigw

import (
	"context"

	"github.com/morzhanov/async-api/api/order"
	"github.com/morzhanov/async-api/internal/mq"
)

type client struct {
	createOrderMQ  mq.MQ
	processOrderMQ mq.MQ
}

type Client interface {
	CreateOrder(ctx context.Context, msg *order.CreateOrderMessage) error
	ProcessOrder(ctx context.Context, orderID string) error
}

func (c *client) CreateOrder(ctx context.Context, msg *order.CreateOrderMessage) error {
	return c.createOrderMQ.WriteMessage(ctx, msg)
}

func (c *client) ProcessOrder(ctx context.Context, orderID string) error {
	return c.processOrderMQ.WriteMessage(ctx, &order.ProcessOrderMessage{ID: orderID})
}

func NewClient(createOrderMQ mq.MQ, processOrderMQ mq.MQ) Client {
	return &client{createOrderMQ, processOrderMQ}
}
