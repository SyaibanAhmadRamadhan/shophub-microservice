package users

import "context"

type RepositoryWriter interface {
	Create(ctx context.Context, input CreateInput) (id int64, err error)
	Update(ctx context.Context, input UpdateInput) error
}

type RepositoryReader interface {
	FindAll(ctx context.Context, input FindAllInput) (output FindAllOutput, err error)
	FindOne(ctx context.Context, input FindOneInput) (output Entity, err error)
}
