package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/morzhanov/async-api/internal/config"
	"github.com/morzhanov/async-api/internal/logger"
	"github.com/morzhanov/async-api/internal/payment"
	"github.com/morzhanov/async-api/internal/psql"
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
	p, err := psql.NewDb(c.PostgresURL)
	failOnError(l, "postgres", err)

	pay := payment.NewPayment(p)
	ctrl, err := payment.NewController(pay, c, l)
	failOnError(l, "service", err)
	go ctrl.Listen(context.Background())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	log.Println("App successfully started!")
	<-quit
	log.Println("received os.Interrupt, exiting...")
}
