package payment

import (
	"context"
	"encoding/json"

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
	c.Meter().IncReqCount()
	et := c.Tracer()("kafka")
	pctx, err := event.GetSpanContext(in)
	if err != nil {
		c.Logger().Error("error during process payment event processing", zap.Error(err))
	}
	sctx, span := et.Start(*pctx, "process-payment")
	defer span.End()

	//res := gpayment.ProcessPaymentMessage{}
	if err := json.Unmarshal(in.Value, &res); err != nil {
		c.Logger().Error("error during process payment event processing", zap.Error(err))
	}
	if err := c.pay.ProcessPayment(sctx, &res); err != nil {
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
	controller, err := event.NewController(c.KafkaURL, c.KafkaTopic, c.KafkaGroupID, log)
	return &eventController{BaseController: controller, pay: pay}, err
}
