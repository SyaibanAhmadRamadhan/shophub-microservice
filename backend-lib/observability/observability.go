package observability

import (
	"fmt"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type OptionParams struct {
	ServiceName    string
	Env            string
	OtlpEndpoint   string
	OtlpUsername   string
	OtlpPassword   string
	KafkaAddr      []string
	TransportKafka *kafka.Transport
}

func NewObservability(params OptionParams) (trace.Tracer, func(), error) {
	closeFunc, err := NewOtel(params.ServiceName, params.OtlpEndpoint, params.OtlpUsername, params.OtlpPassword)
	if err != nil {
		return nil, nil, err
	}

	kafkaTopic := os.Getenv("KAFKA_LOG_TOPIC")
	if kafkaTopic == "" {
		return nil, nil, fmt.Errorf("kafka topic is empty")
	}

	w := &kafka.Writer{
		Addr:            kafka.TCP(params.KafkaAddr...),
		Topic:           kafkaTopic,
		Balancer:        &kafka.LeastBytes{},
		MaxAttempts:     5,
		WriteBackoffMin: time.Duration(100),
		WriteBackoffMax: time.Duration(1 * time.Second),

		BatchSize:    10,
		BatchBytes:   1048576,
		BatchTimeout: time.Duration(3 * time.Second),

		RequiredAcks: kafka.RequireOne,
		Transport:    params.TransportKafka,
	}
	NewLog(w, kafkaTopic, params.Env, params.ServiceName)

	return otel.Tracer(params.ServiceName), func() {
		closeFunc()
		w.Close()
	}, err
}
