package main

import (
	"context"
	"log"
	"order-service/internal/infrastructures"
	"order-service/internal/presentations"
	"order-service/internal/presentations/restapi"
	useraddresses "order-service/internal/repositories/user_addresses"
	"order-service/internal/repositories/users"
	etlservice "order-service/internal/services/etl_service"
	"os"
	"os/signal"
	"syscall"
	"time"

	libkafka "github.com/SyaibanAhmadRamadhan/go-foundation-kit/broker/kafka"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var restApi = &cobra.Command{
	Use:   "rest-api",
	Short: "RUN rest api",
	Run: func(cmd *cobra.Command, args []string) {
		godotenv.Load(".env")

		_, closeObservabilityFn, err := infrastructures.NewObservability()
		if err != nil {
			panic(err)
		}

		_, _, _, closeFn, err := infrastructures.NewPgx()
		if err != nil {
			panic(err)
		}
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		dependency := presentations.Dependency{}

		// inject dependency
		{
			// repository layer
			// userRepository := users.New(rdbms, sq)
			// userAddressesRepository := useraddresses.New(rdbms, sq)

			// service layer
		}

		shutdownServer := restapi.New(dependency)

		<-ctx.Done()
		log.Println("graceful shutdown starting")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		shutdownServer(shutdownCtx)

		done := make(chan struct{})

		go func() {
			closeFn()
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
				log.Println("error: ", err)
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

var rootCmd = &cobra.Command{
	Use:   "mycli",
	Short: "Multi-purpose CLI for app management",
}

func main() {
	rootCmd.AddCommand(restApi, etlConsumer)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
