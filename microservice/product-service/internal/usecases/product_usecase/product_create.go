package productusecase

import (
	"context"
	"fmt"
	"product-service/internal/repositories/products"

	libpgx "github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases/pgx"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/observability"
	"github.com/jackc/pgx/v5"
)

func (p *product) CreateProduct(ctx context.Context, input CreateProductInput) (output CreateProductOutput, err error) {
	err = p.tx.DoTxContext(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite},
		func(ctx context.Context, tx libpgx.RDBMS) (err error) {
			output.ID, err = p.productRepositoryWriter.Create(ctx, products.CreateInput{
				Entity: products.Entity{
					Name:        input.Name,
					Description: input.Description,
					Price:       input.Price,
					Stock:       input.Stock,
					SKU:         input.SKU,
					IsActive:    input.IsActive,
					CategoryID:  input.CategoryID,
					TraceParent: observability.ExtractTraceparent(ctx),
				},
			})
			if err != nil {
				return fmt.Errorf("failed to create product: %w", err)
			}
			return
		},
	)
	if err != nil {
		return output, err
	}

	return
}
