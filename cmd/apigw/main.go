package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/morzhanov/async-api/internal/apigw"
	"github.com/morzhanov/async-api/internal/config"
	"github.com/morzhanov/async-api/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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

	uri := fmt.Sprintf("%s:%s", c.PaymentGRPCurl, c.PaymentGRPCport)
	conn, err := grpc.Dial(uri, grpc.WithInsecure(), grpc.WithBlock())
	failOnError(l, "config", err)
	//client := apigw.NewClient(c.OrderRESTurl, payment.NewPaymentClient(conn))
	srv := apigw.NewController(client, l)
	go srv.Listen(c.APIGWport)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	log.Println("App successfully started!")
	<-quit
	log.Println("received os.Interrupt, exiting...")
}
