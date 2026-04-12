// Package app
//
//nolint:all
package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"abr_paperless_office/internal/handlers"
	"abr_paperless_office/internal/metrics"
	"abr_paperless_office/internal/models"
	adapters "abr_paperless_office/internal/repositories"
	"abr_paperless_office/internal/service"
	appconfig "abr_paperless_office/pkg/config"
	"abr_paperless_office/pkg/database"
	"abr_paperless_office/pkg/logger"
	"abr_paperless_office/pkg/server"

	_ "time/tzdata"

	"github.com/jackc/pgx/v5/pgxpool"
	cnf "github.com/peak/go-config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
)

const ConfigPath = "data/config.yml"

func Run() {
	logger.InitFromEnv()

	ctx := context.Background()
	ch, err := cnf.Watch(ctx, ConfigPath)
	if err != nil {
		log.Printf("failed to watch config: %v\n", err)
		return
	}

	for {
		appconfig.Setup(ConfigPath)
		cfg := appconfig.GetConfig()

		if err := logger.InitEventLogFile(cfg.Log.EventPath); err != nil {
			logger.ErrorCtx(nil, err, "event_log_init_failed")
		}

		pool, err := database.InitDB(ctx, cfg)
		if err != nil {
			log.Println(err)
			_ = logger.CloseEventLogFile()
			return
		}

		rep := adapters.NewRepository(pool)
		services := service.NewService(rep, crmClient)
		handlers := handler.NewHandler(services, cfg)
		srv := server.NewServer(cfg, handlers.Init(), ConfigPath)

		runErrCh := make(chan error, 1)
		go func() {
			if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				runErrCh <- err
				return
			}
			runErrCh <- nil
		}()

		log.Println("🚀 Running on port http://localhost:" + cfg.Server.Port + "/")
		needReload := false
		needExit := false

		for !needReload && !needExit {
			select {
			case err := <-runErrCh:
				if err != nil {
					log.Printf("server failed: %v\n", err)
				}
				needExit = true

			case e, ok := <-ch:
				if !ok {
					log.Printf("config watch channel closed")
					needExit = true
					continue
				}
				if e != nil {
					logger.ErrorCtx(nil, e, "config_watch_error")
					continue
				}

				logger.InfoCtx(nil, "config_changed_reloading")

				ct, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				_ = srv.Stop(ct)
				_ = logger.CloseEventLogFile()
				cancel()

				_ = <-runErrCh
				needReload = true
			}
		}
		pool.Close()
		if needExit {
			_ = logger.CloseEventLogFile()
			return
		}
	}
}
