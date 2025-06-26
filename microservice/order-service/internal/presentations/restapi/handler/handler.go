package handler

import (
	"order-service/internal/presentations"
)

type Handler struct {
	serv *presentations.Dependency
}

type Options struct {
	Serv *presentations.Dependency
}

func NewHandler(opts Options) *Handler {
	return &Handler{
		serv: opts.Serv,
	}
}
