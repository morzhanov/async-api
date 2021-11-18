package order

import (
	"context"
	"encoding/json"

	apiorder "github.com/morzhanov/async-api/api/order"
	"github.com/morzhanov/async-api/api/payment"
	"github.com/morzhanov/async-api/internal/config"
	"github.com/morzhanov/async-api/internal/event"
	"github.com/morzhanov/async-api/internal/mq"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Message struct {
	ID     string
	Name   string
	Amount int
	Status string
}

type createController struct {
	event.BaseController
	coll *mongo.Collection
}

type processController struct {
	event.BaseController
	coll             *mongo.Collection
	processPaymentMq mq.MQ
}

type Controller interface {
	Listen(ctx context.Context)
}

func (c *createController) createOrder(in *kafka.Message) {
	ctx := context.Background()
	msg := apiorder.CreateOrderMessage{}
	if err := json.Unmarshal(in.Value, &msg); err != nil {
		c.Logger().Error("error during create order event processing", zap.Error(err))
	}
	item := Message{ID: uuid.NewV4().String(), Name: msg.Name, Amount: msg.Amount, Status: "new"}
	if _, err := c.coll.InsertOne(ctx, &item); err != nil {
		c.Logger().Error("error during create order event processing", zap.Error(err))
	}
}

func (c *createController) Listen(ctx context.Context) {
	c.BaseController.Listen(ctx, c.createOrder)
}

func (c *processController) processOrder(in *kafka.Message) {
	ctx := context.Background()
	msg := apiorder.ProcessOrderMessage{}
	if err := json.Unmarshal(in.Value, &msg); err != nil {
		c.Logger().Error("error during process order event processing", zap.Error(err))
	}

	filter := bson.D{{"_id", msg.ID}}
	update := bson.D{{"$set", bson.D{{"status", "processed"}}}}
	if _, err := c.coll.UpdateOne(ctx, filter, update); err != nil {
		c.Logger().Error("error during process order event processing", zap.Error(err))
		return
	}
	res := c.coll.FindOne(ctx, filter)
	if res.Err() != nil {
		c.Logger().Error("error during process order event processing", zap.Error(res.Err()))
		return
	}
	var orderMsg Message
	if err := res.Decode(&orderMsg); err != nil {
		c.Logger().Error("error during process order event processing", zap.Error(err))
		return
	}
	if err := c.processPaymentMq.WriteMessage(ctx, &payment.ProcessPaymentMessage{OrderID: orderMsg.ID, Name: orderMsg.Name, Amount: orderMsg.Amount, Status: orderMsg.Status}); err != nil {
		c.Logger().Error("error during process order event processing", zap.Error(err))
		return
	}
}

func (c *processController) Listen(ctx context.Context) {
	c.BaseController.Listen(ctx, c.processOrder)
}

func NewController(
	c *config.Config,
	log *zap.Logger,
	processPaymentMq mq.MQ,
	coll *mongo.Collection,
) ([]Controller, error) {
	createCtrl, err := event.NewController(c.KafkaURL, "order.create", "order.create", log)
	processCtrl, err := event.NewController(c.KafkaURL, "order.process", "order.process", log)
	return []Controller{
		&createController{BaseController: createCtrl, coll: coll},
		&processController{BaseController: processCtrl, coll: coll, processPaymentMq: processPaymentMq},
	}, err
}
