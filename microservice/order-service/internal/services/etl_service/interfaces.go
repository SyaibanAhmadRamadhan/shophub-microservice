package etlservice

import "context"

type UserService interface {
	EtlUsers(ctx context.Context) (err error)
}
