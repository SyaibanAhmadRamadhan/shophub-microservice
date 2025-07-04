package productcategories

import (
	"context"
	"errors"
	"time"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases"
)

func (r *productCategories) Create(ctx context.Context, input CreateInput) (id int64, err error) {
	if input.Tx == nil {
		return id, errors.New("failed to create product category, transaction database is nil")
	}

	err = input.Tx.QueryRow(ctx,
		`INSERT INTO product_categories (name, description, trace_parent) 
		 VALUES ($1, $2, $3) RETURNING id`,
		input.Entity.Name,
		input.Entity.Description,
		input.Entity.TraceParent,
	).Scan(&id)

	if err != nil {
		return id, err
	}

	return
}

func (r *productCategories) Update(ctx context.Context, input UpdateInput) error {
	if input.Tx == nil {
		return errors.New("failed to update product category, transaction database is nil")
	}

	cmdTag, err := input.Tx.Exec(ctx,
		`UPDATE product_categories 
		 SET name = $1, description = $2, trace_parent = $3, updated_at = $4
		 WHERE id = $5`,
		input.Entity.Name,
		input.Entity.Description,
		input.Entity.TraceParent,
		time.Now().UTC(),
		input.Entity.ID,
	)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return databases.ErrNoUpdateRow
	}

	return nil
}
