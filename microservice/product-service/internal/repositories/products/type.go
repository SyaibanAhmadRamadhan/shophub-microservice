package products

import (
	"time"

	libpgx "github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases/pgx"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/utils/primitive"
)

type Entity struct {
	ID          int64     `db:"id"`
	CategoryID  int64     `db:"category_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Price       float64   `db:"price"`
	Stock       int64     `db:"stock"`
	SKU         string    `db:"sku"`
	IsActive    bool      `db:"is_active"`
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
	ID  int64
	SKU string
}

type FindAllInput struct {
	SearchKeyword string
	Pagination    primitive.PaginationInput
}

type FindAllOutput struct {
	Pagination primitive.PaginationOutput
	Entities   []Entity
}
