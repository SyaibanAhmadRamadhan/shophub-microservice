package productusecase

import (
	productcategories "product-service/internal/repositories/product_categories"
	"product-service/internal/repositories/products"

	libpgx "github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases/pgx"
)

type product struct {
	productRepositoryReader         products.RepositoryReader
	productRepositoryWriter         products.RepositoryWriter
	productCategoryRepositoryReader productcategories.RepositoryReader
	productCategoryRepositoryWriter productcategories.RepositoryWriter
	tx                              libpgx.Tx
}

type OptionParams struct {
	ProductRepositoryReader         products.RepositoryReader
	ProductRepositoryWriter         products.RepositoryWriter
	ProductCategoryRepositoryReader productcategories.RepositoryReader
	ProductCategoryRepositoryWriter productcategories.RepositoryWriter
	Tx                              libpgx.Tx
}

func New(optionParams OptionParams) *product {
	return &product{
		productRepositoryReader:         optionParams.ProductRepositoryReader,
		productRepositoryWriter:         optionParams.ProductRepositoryWriter,
		productCategoryRepositoryReader: optionParams.ProductCategoryRepositoryReader,
		productCategoryRepositoryWriter: optionParams.ProductCategoryRepositoryWriter,
		tx:                              optionParams.Tx,
	}
}
