package productusecase

import "context"

type ProductUsecase interface {
	CreateProduct(ctx context.Context, input CreateProductInput) (output CreateProductOutput, err error)
	CreateProductCategory(ctx context.Context, input CreateProductCategoryInput) (output CreateProductCategoryOutput, err error)
}
