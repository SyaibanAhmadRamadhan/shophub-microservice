package logcollector

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log-collector/repositories/provider/loki"
	"time"

	"github.com/segmentio/kafka-go"
)

type logCollector struct {
	kafkaReader *kafka.Reader
	lokiClient  loki.LokiClientInterface
}

func NewLog(kafkaReader *kafka.Reader, lokiClient loki.LokiClientInterface) *logCollector {
	return &logCollector{
		kafkaReader: kafkaReader,
		lokiClient:  lokiClient,
	}
}

func (l *logCollector) Start(ctx context.Context) {
	log.Println("started log collector...")
	for {
		msg, err := l.kafkaReader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				log.Println("Kafka reader canceled by context")
				return
			}
			log.Printf("Error reading message: %v\n", err)
			continue
		}

		select {
		case <-ctx.Done():
			log.Println("Log collector received shutdown signal")
			return
		default:
			log.Printf("Received message: %s\n", string(msg.Value))

			streams := map[string]string{}
			for _, v := range msg.Headers {
				streams[v.Key] = string(v.Value)
			}

			l.lokiClient.MustSendDataLog(ctx, loki.SendDataLogInput{
				Request: loki.SendDataLogRequest{
					Streams: []loki.SendDataLogStramRequest{
						{
							Stream: streams,
							Values: [][2]string{
								{fmt.Sprintf("%d", time.Now().UnixNano()), string(msg.Value)},
							},
						},
					},
				},
			})
		}
	}
}
