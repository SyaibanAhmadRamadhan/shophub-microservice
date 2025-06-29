package users

import (
	"time"

	libpgx "github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases/pgx"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/utils/primitive"
)

type Entity struct {
	ID          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Email       string    `db:"email" json:"email"`
	PhoneNumber string    `db:"phone_number" json:"phoneNumber"`
	Password    string    `db:"password" json:"password"`
	IsVerified  bool      `db:"is_verified" json:"isVerified"`
	TraceParent string    `db:"trace_parent" json:"traceParent"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
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
