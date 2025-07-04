package presentations

import (
	productusecase "product-service/internal/usecases/product_usecase"
)

type Dependency struct {
	ProductUsecase productusecase.ProductUsecase
}
