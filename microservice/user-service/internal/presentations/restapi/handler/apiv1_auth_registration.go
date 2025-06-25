package handler

import (
	"net/http"
	"user-service/.gen/api"
	userservice "user-service/internal/services/user_service"

	libgin "github.com/SyaibanAhmadRamadhan/go-foundation-kit/webserver/gin"
	"github.com/gin-gonic/gin"
)

func (h *Handler) AuthRegistration(c *gin.Context) {
	req := api.AuthRegistrationRequestBody{}

	if ok := libgin.MustShouldBind(c, &req); !ok {
		return
	}

	registerOutput, err := h.serv.UserService.Register(c.Request.Context(), userservice.RegisterInput{
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
	})
	if err != nil {
		libgin.ErrorResponse(c, err)
		return
	}

	resp := api.GeneralResponseSuccessID{
		Id: registerOutput.ID,
	}

	c.JSON(http.StatusCreated, resp)
}
