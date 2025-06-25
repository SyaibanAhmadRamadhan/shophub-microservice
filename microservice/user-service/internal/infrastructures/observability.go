package infrastructures

import (
	"log"
	"os"
	"strings"

	libkafka "github.com/SyaibanAhmadRamadhan/go-foundation-kit/broker/kafka"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/observability"
	"go.opentelemetry.io/otel/trace"
)

func NewObservability() (trace.Tracer, func(), error) {
	observabilityParams := observability.OptionParams{
		ServiceName:  os.Getenv("SERVICE_NAME"),
		Env:          os.Getenv("APP_ENV"),
		OtlpEndpoint: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		OtlpUsername: os.Getenv("OTEL_EXPORTER_OTLP_USERNAME"),
		OtlpPassword: os.Getenv("OTEL_EXPORTER_OTLP_PASSWORD"),
	}
	tracerOtel, closeOtelFn, err := observability.NewObservabilityOtel(observabilityParams)
	if err != nil {
		return nil, nil, err
	}

	kafkaTransport := libkafka.NewTransportSasl(
		os.Getenv("KAFKA_SASL_USER"),
		os.Getenv("KAFKA_SASL_PASS"),
	)
	closeKafkaFn := observability.NewLogWithKafkaHook(
		strings.Split(os.Getenv("KAFKA_ADDRS"), ","),
		kafkaTransport,
		os.Getenv("KAFKA_LOG_TOPIC"),
		observabilityParams,
	)

	log.Println("initialization otel and log hook successfully")
	return tracerOtel, func() {
		closeOtelFn()
		closeKafkaFn()
	}, nil
}
