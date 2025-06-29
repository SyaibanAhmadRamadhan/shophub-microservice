package users

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases"
	"github.com/jackc/pgx/v5"
)

func (r *user) FindOne(ctx context.Context, input FindOneInput) (output Entity, err error) {
	query := r.sq.Select("id", "name", "email", "phone_number", "password", "is_verified").From("users")
	if input.ID != 0 {
		query = query.Where(squirrel.Eq{"id": input.ID})
	}
	if input.Email != "" {
		query = query.Where(squirrel.Eq{"email": input.Email})
	}

	row, err := r.rdbms.QueryRowSq(ctx, query)
	if err != nil {
		return output, err
	}

	err = row.Scan(
		&output.ID,
		&output.Name,
		&output.Email,
		&output.PhoneNumber,
		&output.Password,
		&output.IsVerified,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return output, databases.ErrNoRowFound
		}
		return output, err
	}

	return output, nil
}

func (r *user) FindAll(ctx context.Context, input FindAllInput) (output FindAllOutput, err error) {
	filters := squirrel.And{}

	if input.SearchKeyword != "" {
		filters = append(filters, squirrel.Or{
			squirrel.ILike{"name": "%" + input.SearchKeyword + "%"},
			squirrel.ILike{"email": "%" + input.SearchKeyword + "%"},
			squirrel.ILike{"CAST(id AS TEXT)": "%" + input.SearchKeyword + "%"},
		})
	}

	findQuery := r.sq.
		Select("id", "name", "email", "phone_number", "password", "is_verified").
		From("users")

	countQuery := r.sq.
		Select("COUNT(*)").
		From("users")

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
