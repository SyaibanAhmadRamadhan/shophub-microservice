package useraddresses

import (
	"time"

	libpgx "github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases/pgx"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/utils/primitive"
)

type Entity struct {
	ID          int64     `db:"id"`
	UserID      int64     `db:"user_id"`
	FullAddress int64     `db:"full_address"`
	TraceParent string    `db:"trace_parent"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type UpSertInput struct {
	Entity Entity
	Tx     libpgx.RDBMS
}

type UpdateInput struct {
	Entity Entity
	Tx     libpgx.RDBMS
}

type FindOneInput struct {
	ID int64
}

type FindAllInput struct {
	SearchKeyword string
	Pagination    primitive.PaginationInput
}

type FindAllOutput struct {
	Pagination primitive.PaginationOutput
	Entities   []Entity
}
