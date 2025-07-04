package products

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases"
	"github.com/jackc/pgx/v5"
)

func (r *products) FindOne(ctx context.Context, input FindOneInput) (output Entity, err error) {
	query := r.sq.
		Select("id", "name", "description", "price", "stock", "sku", "is_active", "trace_parent", "created_at", "updated_at").
		From("products")

	if input.ID != 0 {
		query = query.Where(squirrel.Eq{"id": input.ID})
	}
	if input.SKU != "" {
		query = query.Where(squirrel.Eq{"sku": input.SKU})
	}

	row, err := r.rdbms.QueryRowSq(ctx, query)
	if err != nil {
		return output, err
	}

	err = row.Scan(
		&output.ID,
		&output.Name,
		&output.Description,
		&output.Price,
		&output.Stock,
		&output.SKU,
		&output.IsActive,
		&output.TraceParent,
		&output.CreatedAt,
		&output.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return output, databases.ErrNoRowFound
		}
		return output, err
	}

	return output, nil
}
func (r *products) FindAll(ctx context.Context, input FindAllInput) (output FindAllOutput, err error) {
	filters := squirrel.And{}

	if input.SearchKeyword != "" {
		filters = append(filters, squirrel.Or{
			squirrel.ILike{"name": "%" + input.SearchKeyword + "%"},
			squirrel.ILike{"description": "%" + input.SearchKeyword + "%"},
			squirrel.ILike{"sku": "%" + input.SearchKeyword + "%"},
		})
	}

	findQuery := r.sq.
		Select("id", "name", "description", "price", "stock", "sku", "is_active", "trace_parent", "created_at", "updated_at").
		From("products")

	countQuery := r.sq.
		Select("COUNT(*)").
		From("products")

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

	return output, nil
}
