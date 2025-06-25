package users

import (
	"github.com/Masterminds/squirrel"
	libpgx "github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases/pgx"
)

type user struct {
	rdbms libpgx.RDBMS
	sq    squirrel.StatementBuilderType
}

func New(rdbms libpgx.RDBMS, sq squirrel.StatementBuilderType) *user {
	return &user{
		rdbms: rdbms,
		sq:    sq,
	}
}
