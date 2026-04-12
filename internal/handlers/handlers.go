// Package handler
package handler

import (
	"net/http"

	v1 "abr_paperless_office/internal/handlers/v1"
	"abr_paperless_office/internal/middlewares"
	"abr_paperless_office/internal/middlewares/otelgin"
	"abr_paperless_office/internal/service"
	"abr_paperless_office/pkg/config"
	logger "abr_paperless_office/pkg/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/zsais/go-gin-prometheus"
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

func (h *Handler) Init() *gin.Engine {
	rl := middlewares.NewRateLimiter()
	gin.SetMode(h.cnf.Server.Mode)
	gin.DisableConsoleColor()
	router := gin.New()
	router.Use(
		logger.WithRequestID(),
		logger.EventMiddleware(),
		logger.GinMiddleware(),
		gin.Recovery(),
		middlewares.CORS(),
		rl.RateLimitMiddleware(),
		middlewares.APILatencyMiddleware(),
	)
	if h.cnf.Trace.Enabled {
		router.Use(otelgin.Middleware(h.cnf.Server.ServiceName))
	}
	router.NoRoute(middlewares.NoRouteHandler())
	router.NoMethod(middlewares.NoMethodHandler())

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]any{"status": "ok"})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	p := ginprometheus.NewPrometheus("gin")
	ginMetrics := p.HandlerFunc()
	router.Use(func(c *gin.Context) {
		if middlewares.ShouldSkipMetrics(c.Request.URL.Path) {
			c.Next()
			return
		}
		ginMetrics(c)
	})
	p.SetMetricsPath(router)

	h.initAPI(router)
	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services, h.cnf)
	apiPaper := router.Group("/trash/api")
	handlerV1.Init(apiPaper)

}
