package restapi

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"user-service/.gen/api"
	"user-service/internal/presentations"
	"user-service/internal/presentations/restapi/handler"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/validator"
	libgin "github.com/SyaibanAhmadRamadhan/go-foundation-kit/webserver/gin"
)

func New(dependency presentations.Dependency) func(ctx context.Context) {
	validator.InitValidator()
	engine := libgin.NewGin(libgin.GinConfig{
		BlacklistRouteLogResponse: map[string]struct{}{},
		SensitiveFields: map[string]struct{}{
			"password": {},
		},
		Validator: validator.Validate,
		AppName:   os.Getenv("SERVICE_NAME"),
	})
	h := handler.NewHandler(handler.Options{
		Serv: &dependency,
	})
	api.RegisterHandlers(engine, h)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", os.Getenv("APP_PORT")),
		Handler:           engine.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return func(ctx context.Context) {
		if err := srv.Shutdown(ctx); err != nil {
			log.Println("Server Shutdown:", err)
		}

		log.Println("shutdown server successfully")
	}
}
