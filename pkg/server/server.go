// Package server
package server

import (
	"context"
	"net/http"

	"trash_project/pkg/config"
)

type Server struct {
	httpServer *http.Server
	ConfPath   string
}

func NewServer(cfg *config.Configuration, handler http.Handler, cnfpath string) *Server {
	return &Server{
		ConfPath: cnfpath,
		httpServer: &http.Server{
			Addr:           ":" + cfg.Server.Port,
			Handler:        handler,
			ReadTimeout:    cfg.Server.ReadTimeout,
			WriteTimeout:   cfg.Server.WriteTimeout,
			MaxHeaderBytes: cfg.Server.MaxHeaderMegabytes << 20,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
