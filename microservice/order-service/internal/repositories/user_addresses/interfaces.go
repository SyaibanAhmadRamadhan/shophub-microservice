package useraddresses

import (
	"context"

	libpgx "github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases/pgx"
)

type RepositoryWriter interface {
	UpSert(ctx context.Context, input UpSertInput) (id int64, err error)
	Update(ctx context.Context, input UpdateInput) error
	Delete(ctx context.Context, tx libpgx.RDBMS, id int64) error
}

type RepositoryReader interface {
	FindAll(ctx context.Context, input FindAllInput) (output FindAllOutput, err error)
	FindOne(ctx context.Context, input FindOneInput) (output Entity, err error)
}
