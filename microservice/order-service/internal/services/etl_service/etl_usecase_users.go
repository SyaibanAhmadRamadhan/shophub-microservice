package etlservice

import (
	"context"
	"errors"
	"fmt"
	"order-service/internal/repositories/users"
	"os"
	"strings"
	"time"

	libkafka "github.com/SyaibanAhmadRamadhan/go-foundation-kit/broker/kafka"
	libpgx "github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases/pgx"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/observability"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/utils/primitive"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"go.opentelemetry.io/otel/attribute"
)

func (u *etl) EtlUsers(ctx context.Context) (err error) {
	mechanism := plain.Mechanism{
		Username: os.Getenv("KAFKA_SASL_USER"),
		Password: os.Getenv("KAFKA_SASL_PASS"),
	}
	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}
	output, err := u.pubSubKafka.Subscribe(ctx, libkafka.SubInput{
		Config: kafka.ReaderConfig{
			Brokers: strings.Split(os.Getenv("KAFKA_ADDRS"), ","),
			GroupID: os.Getenv("KAKFA_ETL_CONSUMER_GROUP"),
			Topic:   os.Getenv("KAFKA_ETL_USER"),
			Dialer:  dialer,
		},
	})
	if err != nil {
		return err
	}

	for {
		data := primitive.DebeziumExtractNewRecordState[EtlUserEntity]{}
		msg, err := output.Reader.FetchMessage(ctx, &data)
		if err != nil {
			if !errors.Is(err, libkafka.ErrJsonUnmarshal) {
				return err
			}
			continue
		}

		msgCarrier := libkafka.NewMsgCarrier(&msg)
		newCtx := u.propagaion.Extract(context.Background(), msgCarrier)

		newCtx, span := u.tracer.Start(newCtx, "debezium.message.info")
		span.SetAttributes(
			attribute.String("debezium.operation", data.Payload.Op),
			attribute.String("debezium.schema", data.Schema.Name),
			attribute.String("kafka.topic", msg.Topic),
			attribute.String("kafka.partition", fmt.Sprintf("%d", msg.Partition)),
			attribute.Int64("kafka.offset", msg.Offset),
			attribute.String("debezium.source.table", data.Payload.Table),
		)

		err = u.tx.DoTxContext(newCtx, pgx.TxOptions{
			IsoLevel:   pgx.ReadCommitted,
			AccessMode: pgx.ReadWrite,
		}, func(ctx context.Context, tx libpgx.RDBMS) error {
			switch data.Payload.Op {
			case "c", "u":
				_, err = u.userRepositoryWriter.UpSert(ctx, users.UpSertInput{
					Tx:     tx,
					Entity: data.Payload.Entity,
				})
			case "d":
				err = u.userRepositoryWriter.Delete(ctx, tx, data.Payload.ID)
			default:
				log.Warn().Msgf("unsupported operation %s", data.Payload.Op)
			}
			return err
		})
		if err != nil {
			observability.Start(ctx, zerolog.ErrorLevel).Err(err).Msgf("failed %s operation", data.Payload.Op)
		} else {
			err = output.Reader.CommitMessages(newCtx, msg)
			if err != nil {
				observability.Start(ctx, zerolog.ErrorLevel).Err(err).Msgf("failed commit %s operation", data.Payload.Op)
			}
		}
		span.End()
	}
}
