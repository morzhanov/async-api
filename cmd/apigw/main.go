package main

import (
	"log"
	"os"
	"os/signal"

	api "github.com/morzhanov/async-api/api/apigw"
	"github.com/morzhanov/async-api/internal/apigw"
	"github.com/morzhanov/async-api/internal/config"
	"github.com/morzhanov/async-api/internal/logger"
	"github.com/morzhanov/async-api/internal/mq"
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

	err = api.Build(c.KafkaURL, "kafka", c.ProtocolVersion)
	failOnError(l, "api", err)
	createOrderMq, err := mq.NewMq(c.KafkaURL, "order.create")
	failOnError(l, "message_queue", err)
	processOrderMq, err := mq.NewMq(c.KafkaURL, "order.process")
	failOnError(l, "message_queue", err)

	client := apigw.NewClient(createOrderMq, processOrderMq)
	srv := apigw.NewController(client, l)
	go srv.Listen(c.APIGWport)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	log.Println("App successfully started!")
	<-quit
	log.Println("received os.Interrupt, exiting...")
}
