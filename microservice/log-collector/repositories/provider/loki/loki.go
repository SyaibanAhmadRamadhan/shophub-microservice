package loki

import (
	"context"
	"log"
	"os"
	"time"

	"resty.dev/v3"
)

type lokiClient struct {
	resty *resty.Client
}

func NewLokiClient() (*lokiClient, func()) {
	client := resty.New()
	client.AddRetryConditions(func(response *resty.Response, err error) bool {
		return response.StatusCode() >= 500
	})
	return &lokiClient{
			resty: client,
		}, func() {
			err := client.Close()
			if err != nil {
				log.Println("error closing resty client", err)
			}
		}
}

func (l *lokiClient) MustSendDataLog(ctx context.Context, input SendDataLogInput) {
	_, err := l.resty.R().
		SetContext(ctx).
		SetBasicAuth(
			os.Getenv("LOKI_USER"),
			os.Getenv("LOKI_PASS"),
		).
		SetRetryCount(3).
		SetRetryWaitTime(2*time.Second).
		SetRetryMaxWaitTime(10*time.Second).
		SetHeader("Content-Type", "application/json").
		SetBody(input.Request).
		Post(os.Getenv("LOKI_ENDPOINT"))
	if err != nil {
		log.Println("error send data log to loki", err)
	}
}
