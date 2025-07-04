package products

import (
	"context"
	"errors"
	"time"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases"
)

func (r *products) Create(ctx context.Context, input CreateInput) (id int64, err error) {
	if input.Tx == nil {
		return id, errors.New("failed to create product, transaction database is nil")
	}

	err = input.Tx.QueryRow(ctx,
		`INSERT INTO products (
			name, description, price, stock, sku, is_active, trace_parent
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		input.Entity.Name,
		input.Entity.Description,
		input.Entity.Price,
		input.Entity.Stock,
		input.Entity.SKU,
		input.Entity.IsActive,
		input.Entity.TraceParent,
	).Scan(&id)

	if err != nil {
		return id, err
	}

	return
}
func (r *products) Update(ctx context.Context, input UpdateInput) error {
	if input.Tx == nil {
		return errors.New("failed to update product, transaction database is nil")
	}

	cmdTag, err := input.Tx.Exec(ctx,
		`UPDATE products 
		 SET 
			name = $1,
			description = $2,
			price = $3,
			stock = $4,
			sku = $5,
			is_active = $6,
			trace_parent = $7,
			updated_at = $8
		 WHERE id = $9`,
		input.Entity.Name,
		input.Entity.Description,
		input.Entity.Price,
		input.Entity.Stock,
		input.Entity.SKU,
		input.Entity.IsActive,
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
