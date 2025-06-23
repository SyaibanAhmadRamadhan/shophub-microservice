package loki

import "context"

type LokiClientInterface interface {
	MustSendDataLog(ctx context.Context, input SendDataLogInput)
}
