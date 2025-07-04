package productcategories

import (
	"github.com/Masterminds/squirrel"
	libpgx "github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases/pgx"
)

type productCategories struct {
	rdbms libpgx.RDBMS
	sq    squirrel.StatementBuilderType
}

func New(rdbms libpgx.RDBMS, sq squirrel.StatementBuilderType) *productCategories {
	return &productCategories{
		rdbms: rdbms,
		sq:    sq,
	}
}
