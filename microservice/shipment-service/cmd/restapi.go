package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"shipment-service/internal/infrastructures"
	"shipment-service/internal/presentations"
	"shipment-service/internal/presentations/restapi"
	"syscall"
	"time"

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
