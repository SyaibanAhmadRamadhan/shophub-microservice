package main

import (
	"context"
	"fmt"
	"log"
	"log-collector/repositories/provider/loki"
	logcollector "log-collector/services/log-collector"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

func main() {
	godotenv.Load(".env")
	mechanism := plain.Mechanism{
		Username: os.Getenv("KAFKA_SASL_USER"),
		Password: os.Getenv("KAFKA_SASL_PASS"),
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}

	fmt.Println(os.Getenv("KAFKA_BROKER"))

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKER")},
		GroupID:  "consumer-group-log-collector-1",
		Topic:    os.Getenv("KAFKA_LOG_TOPIC"),
		Dialer:   dialer,
		MaxBytes: 10e6,
	})

	lokiClient, closeFn := loki.NewLokiClient()

	logCollectorService := logcollector.NewLog(r, lokiClient)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logCollectorService.Start(ctx)
	}()

	<-ctx.Done()
	log.Println("shutting down gracefully...")
	closeFn()
	r.Close()
}
