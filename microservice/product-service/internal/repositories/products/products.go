package products

import (
	"github.com/Masterminds/squirrel"
	libpgx "github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases/pgx"
)

type products struct {
	rdbms libpgx.RDBMS
	sq    squirrel.StatementBuilderType
}

func New(rdbms libpgx.RDBMS, sq squirrel.StatementBuilderType) *products {
	return &products{
		rdbms: rdbms,
		sq:    sq,
	}
}
