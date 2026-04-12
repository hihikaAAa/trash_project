// Package v1
package v1

import (
	"abr_paperless_office/internal/service"
	"abr_paperless_office/pkg/config"

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
		links := v1.Group("/links")
		{
			h.initLinksRoutes(links)
		}
		docs := v1.Group("/docs")
		{
			h.initDocsRoutes(docs)
		}
		logs := v1.Group("/logs")
		{
			h.initLogsRoutes(logs)
		}
	}
}
