package userservice

import "context"

type UserService interface {
	Register(ctx context.Context, input RegisterInput) (output RegisterOutput, err error)
}
