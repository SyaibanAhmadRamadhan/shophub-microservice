package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-service/internal/infrastructures"
	"user-service/internal/presentations"
	"user-service/internal/presentations/restapi"
	useraddresses "user-service/internal/repositories/user_addresses"
	"user-service/internal/repositories/users"
	userservice "user-service/internal/services/user_service"

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

		rdbms, tx, sq, closeFn, err := infrastructures.NewPgx()
		if err != nil {
			panic(err)
		}
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		dependency := presentations.Dependency{}

		// inject dependency
		{
			// repository layer
			userRepository := users.New(rdbms, sq)
			userAddressesRepository := useraddresses.New(rdbms, sq)

			// service layer
			dependency.UserService = userservice.New(userservice.OptionParams{
				UserRepositoryReader:        userRepository,
				UserRepositoryWriter:        userRepository,
				UserAddressRepositoryWriter: userAddressesRepository,
				UserAddressRepositoryReader: userAddressesRepository,
				Tx:                          tx,
			})
		}

		shutdownServer := restapi.New(dependency)

		<-ctx.Done()
		slog.Info("graceful shutdown starting")

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
			slog.Info("graceful shutdown completed")
		case <-shutdownCtx.Done():
			slog.Warn("graceful shutdown timed out", slog.Any("error", shutdownCtx.Err()))
		}
	},
}

var rootCmd = &cobra.Command{
	Use:   "mycli",
	Short: "Multi-purpose CLI for app management",
}

func main() {
	rootCmd.AddCommand(restApi)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
