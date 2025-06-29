package etlservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	useraddresses "order-service/internal/repositories/user_addresses"
	"order-service/internal/repositories/users"
	"os"
	"strings"

	libkafka "github.com/SyaibanAhmadRamadhan/go-foundation-kit/broker/kafka"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases"
	libpgx "github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases/pgx"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/observability"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/utils/primitive"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/attribute"
)

func (u *etl) EtlUserAddress(ctx context.Context) (err error) {
	output, err := u.pubSubKafka.Subscribe(ctx, libkafka.SubInput{
		Config: kafka.ReaderConfig{
			Brokers: strings.Split(os.Getenv("KAFKA_ADDRS"), ","),
			GroupID: os.Getenv("KAKFA_ETL_USER_ADDRESS_CONSUMER_GROUP"),
			Topic:   os.Getenv("KAFKA_ETL_USER_ADDRESS"),
			Dialer:  u.kafkaDialer,
		},
	})
	if err != nil {
		return err
	}

	for {
		data := primitive.DebeziumExtractNewRecordState[EtlUserAddressEntity]{}
		msg, err := output.Reader.FetchMessage(ctx, &data)
		if err != nil {
			if !errors.Is(err, libkafka.ErrJsonUnmarshal) {
				return err
			}
			continue
		}
		select {
		case <-ctx.Done():
			slog.Info("order service received shutdown signal")
			return nil
		default:
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
					userData, errFindOneUser := u.userRepositoryReader.FindOne(ctx, users.FindOneInput{
						ID: data.Payload.UserID,
					})
					if errFindOneUser != nil {
						if !errors.Is(errFindOneUser, databases.ErrNoRowFound) {
							return errFindOneUser
						}
					}
					span.SetAttributes(attribute.Bool("user.existing", userData.ID > 0))

					_, err = u.userAddressRepositoryWriter.UpSert(ctx, useraddresses.UpSertInput{
						Tx:     tx,
						Entity: data.Payload.Entity,
					})
				case "d":
					err = u.userAddressRepositoryWriter.Delete(ctx, tx, data.Payload.ID)
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
}
