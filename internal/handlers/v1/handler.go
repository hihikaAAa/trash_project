// Package v1 contains HTTP handlers for API v1.
package v1

import (
	"github.com/hihikaAAa/trash_project/internal/service"
	"github.com/hihikaAAa/trash_project/pkg/config"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
	cnf      *config.Configuration
}

func NewHandler(services *service.Service, cnf *config.Configuration) *Handler {
	return &Handler{
		services: services,
		cnf:      cnf,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		orders := v1.Group("/orders")
		{
			h.initOrdersRoutes(orders)
		}
	}
}
