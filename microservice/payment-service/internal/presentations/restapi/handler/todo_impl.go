package handler

import "github.com/gin-gonic/gin"

func (h *Handler) AuthLogin(c *gin.Context)    {}
func (h *Handler) RefreshToken(c *gin.Context) {}

func (h *Handler) AuthorizationAuth(c *gin.Context) {}
func (h *Handler) AuthRegistration(c *gin.Context)  {}
