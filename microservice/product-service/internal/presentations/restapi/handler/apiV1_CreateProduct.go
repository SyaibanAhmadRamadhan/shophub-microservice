package handler

import (
	"net/http"
	"product-service/.gen/api"
	productusecase "product-service/internal/usecases/product_usecase"

	libgin "github.com/SyaibanAhmadRamadhan/go-foundation-kit/webserver/gin"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateProduct(c *gin.Context) {
	req := api.CreateProductJSONRequestBody{}

	if ok := libgin.MustShouldBind(c, &req); !ok {
		return
	}

	createProductOutput, err := h.serv.ProductUsecase.CreateProduct(c.Request.Context(), productusecase.CreateProductInput{
		CategoryID:  req.CategoryId,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		SKU:         req.Sku,
		IsActive:    req.IsActive,
	})
	if err != nil {
		libgin.ErrorResponse(c, err)
		return
	}

	resp := api.GeneralResponseSuccessID{
		Id: createProductOutput.ID,
	}

	c.JSON(http.StatusCreated, resp)
}
