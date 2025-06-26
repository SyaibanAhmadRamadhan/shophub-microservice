package useraddresses

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases"
	"github.com/jackc/pgx/v5"
)

func (r *userAddresses) FindOne(ctx context.Context, input FindOneInput) (output Entity, err error) {
	err = r.rdbms.QueryRow(ctx,
		`SELECT id, user_id, full_address
		 FROM user_addresses
		 WHERE id = $1`,
		input.ID,
	).Scan(
		&output.ID,
		&output.UserID,
		&output.FullAddress,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return output, databases.ErrNoRowFound
		}
		return output, err
	}

	return output, nil
}

func (r *userAddresses) FindAll(ctx context.Context, input FindAllInput) (output FindAllOutput, err error) {
	filters := squirrel.And{}

	if input.SearchKeyword != "" {
		filters = append(filters, squirrel.Or{
			squirrel.ILike{"full_address": "%" + input.SearchKeyword + "%"},
			squirrel.ILike{"CAST(user_id AS TEXT)": "%" + input.SearchKeyword + "%"},
		})
	}

	findQuery := r.sq.
		Select("user_id", "full_address", "id").
		From("user_addresses")

	countQuery := r.sq.
		Select("COUNT(*)").
		From("user_addresses")

	if len(filters) > 0 {
		countQuery = countQuery.Where(filters)
		findQuery = findQuery.Where(filters)
	}

	rows, paginationOutput, err := r.rdbms.QuerySqPagination(ctx, countQuery, findQuery, input.Pagination)
	if err != nil {
		return output, err
	}
	defer rows.Close()

	output.Entities, err = pgx.CollectRows(rows, pgx.RowToStructByName[Entity])
	if err != nil {
		return output, err
	}

	output.Pagination = paginationOutput

	return
}
