package productusecase

type CreateProductInput struct {
	CategoryID  int64
	Name        string
	Description string
	Price       float64
	Stock       int64
	SKU         string
	IsActive    bool
}

type CreateProductOutput struct {
	ID int64
}

type CreateProductCategoryInput struct {
	Name        string
	Description string
}

type CreateProductCategoryOutput struct {
	ID int64
}
