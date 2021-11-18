package payment

import (
	"context"
	"encoding/json"

	apipayment "github.com/morzhanov/async-api/api/payment"
	"github.com/morzhanov/async-api/internal/config"
	"github.com/morzhanov/async-api/internal/event"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type eventController struct {
	event.BaseController
	pay Payment
}

type Controller interface {
	Listen(ctx context.Context)
}

func (c *eventController) processPayment(in *kafka.Message) {
	ctx := context.Background()
	res := apipayment.ProcessPaymentMessage{}
	if err := json.Unmarshal(in.Value, &res); err != nil {
		c.Logger().Error("error during process payment event processing", zap.Error(err))
	}
	if err := c.pay.ProcessPayment(ctx, &res); err != nil {
		c.Logger().Error("error during process payment event processing", zap.Error(err))
	}
}

func (c *eventController) Listen(ctx context.Context) {
	c.BaseController.Listen(ctx, c.processPayment)
}

func NewController(
	pay Payment,
	c *config.Config,
	log *zap.Logger,
) (Controller, error) {
	controller, err := event.NewController(c.KafkaURL, "payment.process", "payment.process", log)
	return &eventController{BaseController: controller, pay: pay}, err
}
