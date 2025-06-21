package observability

import (
	"context"
	"encoding/json"
	"fmt"

	stdlog "log"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/trace"
)

func NewLog(kafka *kafka.Writer, topic, env, serviceName string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(&kafkaHook{
		writer:      kafka,
		topic:       topic,
		env:         env,
		serviceName: serviceName,
	}).With().
		Str("env", env).
		Str("service_name", serviceName).
		Timestamp().Logger()
	log.Logger = logger
}

type kafkaHook struct {
	writer      *kafka.Writer
	topic       string
	env         string
	serviceName string
}

func (w *kafkaHook) Write(p []byte) (n int, err error) {
	if w.env != "production" {
		stdlog.Printf(string(p))
	}

	var payload map[string]any
	if err := json.Unmarshal(p, &payload); err != nil {
		stdlog.Printf("KafkaLogWriter: Failed to parse log JSON:", err)
		return len(p), nil
	}

	level := payload["level"]
	statusCode := payload["status_code"]
	spanID := payload["span_id"]
	traceID := payload["trace_id"]
	headers := []kafka.Header{
		{Key: "service_name", Value: []byte(w.serviceName)},
		{Key: "env", Value: []byte(w.env)},
		{Key: "level", Value: []byte(fmt.Sprintf("%v", level))},
	}

	if statusCode != nil {
		headers = append(headers, kafka.Header{
			Key: "status_code", Value: []byte(fmt.Sprintf("%v", statusCode)),
		})
	}
	if spanID != nil {
		headers = append(headers, kafka.Header{
			Key: "span_id", Value: []byte(fmt.Sprintf("%v", spanID)),
		})
	}
	if traceID != nil {
		headers = append(headers, kafka.Header{
			Key: "trace_id", Value: []byte(fmt.Sprintf("%v", traceID)),
		})
	}

	if traceID == nil {
		stdlog.Printf("warn: invalid log, must be need trace_id")
		return
	}

	err = w.writer.WriteMessages(context.Background(), kafka.Message{
		Value:   p,
		Headers: headers,
		Key:     []byte(fmt.Sprintf("%v", traceID)),
		Topic:   w.topic,
	})
	if err != nil {
		stdlog.Printf("KafkaLogWriter: Failed to send log to Kafka:", err)
	}

	return len(p), nil
}

func Start(ctx context.Context, level zerolog.Level) *zerolog.Event {
	traceID := ""
	spanID := ""
	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.IsValid() {
		traceID = spanContext.TraceID().String()
		spanID = spanContext.SpanID().String()
	}

	switch level {
	case zerolog.TraceLevel:
		return log.Trace().Str("trace_id", traceID).Str("span_id", spanID)
	case zerolog.DebugLevel:
		return log.Debug().Str("trace_id", traceID).Str("span_id", spanID)
	case zerolog.InfoLevel:
		return log.Info().Str("trace_id", traceID).Str("span_id", spanID)
	case zerolog.WarnLevel:
		return log.Warn().Str("trace_id", traceID).Str("span_id", spanID)
	case zerolog.ErrorLevel:
		return log.Error().Str("trace_id", traceID).Str("span_id", spanID)
	case zerolog.FatalLevel:
		return log.Fatal().Str("trace_id", traceID).Str("span_id", spanID)
	case zerolog.PanicLevel:
		return log.Panic().Str("trace_id", traceID).Str("span_id", spanID)
	default:
		return log.Info().Str("trace_id", traceID).Str("span_id", spanID)
	}
}
