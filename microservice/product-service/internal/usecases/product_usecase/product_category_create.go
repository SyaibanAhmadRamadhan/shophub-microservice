package productusecase

import (
	"context"
	"fmt"
	productcategories "product-service/internal/repositories/product_categories"

	libpgx "github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases/pgx"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/observability"
	"github.com/jackc/pgx/v5"
)

func (p *product) CreateProductCategory(ctx context.Context, input CreateProductCategoryInput) (output CreateProductCategoryOutput, err error) {
	err = p.tx.DoTxContext(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite},
		func(ctx context.Context, tx libpgx.RDBMS) (err error) {
			output.ID, err = p.productCategoryRepositoryWriter.Create(ctx, productcategories.CreateInput{
				Entity: productcategories.Entity{
					Name:        input.Name,
					Description: input.Description,
					TraceParent: observability.ExtractTraceparent(ctx),
				},
			})
			if err != nil {
				return fmt.Errorf("failed create product categories: %w", err)
			}
			return
		},
	)

	return
}
