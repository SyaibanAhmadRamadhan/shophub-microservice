package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"shipment-service/internal/infrastructures"
	useraddresses "shipment-service/internal/repositories/user_addresses"
	"shipment-service/internal/repositories/users"
	etlservice "shipment-service/internal/services/etl_service"
	"syscall"
	"time"

	libkafka "github.com/SyaibanAhmadRamadhan/go-foundation-kit/broker/kafka"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var etlConsumer = &cobra.Command{
	Use:   "etl",
	Short: "RUN ETL",
	Run: func(cmd *cobra.Command, args []string) {
		godotenv.Load(".env")

		tracerOtel, closeObservabilityFn, err := infrastructures.NewObservability()
		if err != nil {
			panic(err)
		}

		rdbms, tx, sq, closeFn, err := infrastructures.NewPgx()
		if err != nil {
			panic(err)
		}

		pubSubKafka := libkafka.New(libkafka.WithOtel())

		// repository layer
		userRepository := users.New(rdbms, sq)
		userAddressesRepository := useraddresses.New(rdbms, sq)

		// service layer
		etlService := etlservice.New(etlservice.OptionParams{
			UserRepositoryReader:        userRepository,
			UserRepositoryWriter:        userRepository,
			UserAddressRepositoryWriter: userAddressesRepository,
			UserAddressRepositoryReader: userAddressesRepository,
			Tx:                          tx,
			PubSubKafka:                 pubSubKafka,
			Tracer:                      tracerOtel,
		})

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		go func() {
			err := etlService.EtlUsers(ctx)
			if err != nil {
				log.Println("error et user: ", err)
				stop()
			}
			err = etlService.EtlUserAddress(ctx)
			if err != nil {
				log.Println("error etl user address: ", err)
				stop()
			}
		}()

		<-ctx.Done()
		log.Println("graceful shutdown starting")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		done := make(chan struct{})

		go func() {
			closeFn()
			pubSubKafka.Close()
			closeObservabilityFn()
			close(done)
		}()

		select {
		case <-done:
			log.Println("graceful shutdown completed")
		case <-shutdownCtx.Done():
			log.Println("graceful shutdown timed out")
		}
	},
}
