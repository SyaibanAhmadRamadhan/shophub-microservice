package etlservice

import (
	"os"
	useraddresses "shipment-service/internal/repositories/user_addresses"
	"shipment-service/internal/repositories/users"
	"time"

	libkafka "github.com/SyaibanAhmadRamadhan/go-foundation-kit/broker/kafka"
	libpgx "github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases/pgx"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type etl struct {
	userRepositoryReader        users.RepositoryReader
	userRepositoryWriter        users.RepositoryWriter
	userAddressRepositoryWriter useraddresses.RepositoryWriter
	userAddressRepositoryReader useraddresses.RepositoryReader
	tx                          libpgx.Tx
	pubSubKafka                 libkafka.PubSub

	propagaion  propagation.TextMapPropagator
	tracer      trace.Tracer
	kafkaDialer *kafka.Dialer
}

type OptionParams struct {
	UserRepositoryReader        users.RepositoryReader
	UserRepositoryWriter        users.RepositoryWriter
	UserAddressRepositoryWriter useraddresses.RepositoryWriter
	UserAddressRepositoryReader useraddresses.RepositoryReader
	PubSubKafka                 libkafka.PubSub
	Tx                          libpgx.Tx

	Tracer trace.Tracer
}

func New(optionParams OptionParams) *etl {
	mechanism := plain.Mechanism{
		Username: os.Getenv("KAFKA_SASL_USER"),
		Password: os.Getenv("KAFKA_SASL_PASS"),
	}
	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}

	return &etl{
		userRepositoryReader:        optionParams.UserRepositoryReader,
		userRepositoryWriter:        optionParams.UserRepositoryWriter,
		userAddressRepositoryWriter: optionParams.UserAddressRepositoryWriter,
		userAddressRepositoryReader: optionParams.UserAddressRepositoryReader,
		tx:                          optionParams.Tx,
		pubSubKafka:                 optionParams.PubSubKafka,
		tracer:                      optionParams.Tracer,
		propagaion:                  otel.GetTextMapPropagator(),
		kafkaDialer:                 dialer,
	}
}
