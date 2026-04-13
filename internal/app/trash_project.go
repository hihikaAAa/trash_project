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

	handler "github.com/hihikaAAa/trash_project/internal/handlers"
	adapters "github.com/hihikaAAa/trash_project/internal/repositories"
	"github.com/hihikaAAa/trash_project/internal/service"
	appconfig "github.com/hihikaAAa/trash_project/pkg/config"
	"github.com/hihikaAAa/trash_project/pkg/database"
	"github.com/hihikaAAa/trash_project/pkg/logger"
	"github.com/hihikaAAa/trash_project/pkg/server"

	_ "time/tzdata"

	cnf "github.com/peak/go-config"
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
		services := service.NewService(rep)
		h := handler.NewHandler(services, cfg)
		srv := server.NewServer(cfg, h.Init(), ConfigPath)

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
