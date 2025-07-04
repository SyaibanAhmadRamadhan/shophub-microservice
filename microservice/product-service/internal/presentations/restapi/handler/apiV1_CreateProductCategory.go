package handler

import (
	"net/http"
	"product-service/.gen/api"
	productusecase "product-service/internal/usecases/product_usecase"

	libgin "github.com/SyaibanAhmadRamadhan/go-foundation-kit/webserver/gin"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateProductCategory(c *gin.Context) {
	req := api.CreateProductCategoryJSONRequestBody{}

	if ok := libgin.MustShouldBind(c, &req); !ok {
		return
	}

	createProductCategoryOutput, err := h.serv.ProductUsecase.CreateProductCategory(c.Request.Context(), productusecase.CreateProductCategoryInput{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		libgin.ErrorResponse(c, err)
		return
	}

	resp := api.GeneralResponseSuccessID{
		Id: createProductCategoryOutput.ID,
	}

	c.JSON(http.StatusCreated, resp)
}
