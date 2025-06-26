package logcollector

import (
	"context"
	"errors"
	"fmt"
	"log-collector/repositories/provider/loki"
	"log/slog"
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
	slog.Info("started log collector...")
	for {
		msg, err := l.kafkaReader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				slog.Warn("kafka reader canceled by context",
					slog.String("reason", err.Error()),
				)
				return
			}
			slog.Error("error reading kafka message",
				slog.Any("error", err),
			)
			continue
		}

		select {
		case <-ctx.Done():
			slog.Info("log collector received shutdown signal")
			return

		default:
			slog.Info("received kafka message",
				slog.String("value", string(msg.Value)),
			)

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
								{
									fmt.Sprintf("%d", time.Now().UnixNano()),
									string(msg.Value),
								},
							},
						},
					},
				},
			})
		}
	}
}
