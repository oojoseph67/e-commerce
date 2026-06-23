package handlers

import "github.com/oojoseph67/ecommerce/internal/services"

type Handler struct {
	service *services.Services
}

func NewHandler(service *services.Services) *Handler {
	return &Handler{service: service}
}
