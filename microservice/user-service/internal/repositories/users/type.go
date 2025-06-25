package users

import (
	"time"

	libpgx "github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases/pgx"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/utils/primitive"
)

type Entity struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Email       string    `db:"email"`
	PhoneNumber string    `db:"phone_number"`
	Password    string    `db:"password"`
	IsVerified  bool      `db:"is_verified"`
	TraceParent string    `db:"trace_parent"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type CreateInput struct {
	Entity Entity
	Tx     libpgx.RDBMS
}

type UpdateInput struct {
	Entity Entity
	Tx     libpgx.RDBMS
}

type FindOneInput struct {
	ID    int64
	Email string
}

type FindAllInput struct {
	SearchKeyword string
	Pagination    primitive.PaginationInput
}

type FindAllOutput struct {
	Pagination primitive.PaginationOutput
	Entities   []Entity
}
