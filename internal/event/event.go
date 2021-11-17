package event

import (
	"context"
	"encoding/json"

	"github.com/morzhanov/async-api/internal/mq"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type baseController struct {
	mq      mq.MQ
	log     *zap.Logger
	groupID string
}

type BaseController interface {
	Listen(ctx context.Context, processRequest func(*kafka.Message))
	ConsumerGroupId() string
	Logger() *zap.Logger
}

func (c *baseController) Listen(
	ctx context.Context,
	processRequest func(*kafka.Message),
) {
	r := c.mq.CreateReader(c.groupID)
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			c.log.Error(err.Error())
			continue
		}
		go processRequest(&m)
		select {
		case <-ctx.Done():
			break
		default:
			continue
		}
	}
}

func (c *baseController) Logger() *zap.Logger     { return c.log }
func (c *baseController) ConsumerGroupId() string { return c.groupID }

func GetSpanContext(msg *kafka.Message) (*context.Context, error) {
	var h kafka.Header
	for _, v := range msg.Headers {
		if v.Key == "span-context" {
			h = v
			break
		}
	}
	var sctx context.Context
	if err := json.Unmarshal(h.Value, &sctx); err != nil {
		return nil, err
	}
	return &sctx, nil
}

func NewController(
	kafkaUrl string,
	kafkaTopic string,
	kafkaGroupID string,
	log *zap.Logger,
) (BaseController, error) {
	msgQ, err := mq.NewMq(kafkaUrl, kafkaTopic)
	return &baseController{mq: msgQ, groupID: kafkaGroupID, log: log}, err
}
