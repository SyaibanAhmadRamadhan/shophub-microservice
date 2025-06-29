package useraddresses

import (
	"github.com/Masterminds/squirrel"
	libpgx "github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases/pgx"
)

type userAddresses struct {
	rdbms libpgx.RDBMS
	sq    squirrel.StatementBuilderType
}

func New(rdbms libpgx.RDBMS, sq squirrel.StatementBuilderType) *userAddresses {
	return &userAddresses{
		rdbms: rdbms,
		sq:    sq,
	}
}
