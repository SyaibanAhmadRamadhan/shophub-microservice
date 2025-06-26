package etlservice

import (
	useraddresses "order-service/internal/repositories/user_addresses"
	"order-service/internal/repositories/users"

	libkafka "github.com/SyaibanAhmadRamadhan/go-foundation-kit/broker/kafka"
	libpgx "github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases/pgx"
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

	propagaion propagation.TextMapPropagator
	tracer     trace.Tracer
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
	return &etl{
		userRepositoryReader:        optionParams.UserRepositoryReader,
		userRepositoryWriter:        optionParams.UserRepositoryWriter,
		userAddressRepositoryWriter: optionParams.UserAddressRepositoryWriter,
		userAddressRepositoryReader: optionParams.UserAddressRepositoryReader,
		tx:                          optionParams.Tx,
		pubSubKafka:                 optionParams.PubSubKafka,
		tracer:                      optionParams.Tracer,
		propagaion:                  otel.GetTextMapPropagator(),
	}
}
