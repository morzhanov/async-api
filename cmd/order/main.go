package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	api "github.com/morzhanov/async-api/api/order"
	"github.com/morzhanov/async-api/internal/config"
	"github.com/morzhanov/async-api/internal/logger"
	"github.com/morzhanov/async-api/internal/mongodb"
	"github.com/morzhanov/async-api/internal/mq"
	"github.com/morzhanov/async-api/internal/order"
	"go.uber.org/zap"
)

func failOnError(l *zap.Logger, step string, err error) {
	if err != nil {
		l.Fatal("initialization error", zap.Error(err), zap.String("step", step))
	}
}

func main() {
	l, err := logger.NewLogger()
	if err != nil {
		log.Fatal("initialization error during logger setup")
	}
	c, err := config.NewConfig()
	failOnError(l, "config", err)
	m, err := mongodb.NewMongoDB(c.MongoURL)
	failOnError(l, "mongodb", err)

	err = api.Build(c.KafkaURL, "kafka", c.ProtocolVersion)
	failOnError(l, "api", err)
	processPaymentMq, err := mq.NewMq(c.KafkaURL, "payment.process")
	failOnError(l, "message_queue", err)

	ctrls, err := order.NewController(c, l, processPaymentMq, m)
	failOnError(l, "message_queue", err)
	for _, ctrl := range ctrls {
		ctx := context.Background()
		go ctrl.Listen(ctx)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	log.Println("App successfully started!")
	<-quit
	log.Println("received os.Interrupt, exiting...")
}
