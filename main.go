package main

import (
	"context"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/f4nt0md3v/gymshark-packs/config"
	"github.com/f4nt0md3v/gymshark-packs/server/api"
	"github.com/f4nt0md3v/gymshark-packs/server/api/controllers"
	"github.com/f4nt0md3v/gymshark-packs/server/api/services"
)

// @title PACKER REST API
// @version v0.0.1
// @description Packer REST API.

// @BasePath /api
// @query.collection.format multi
func main() {
	conf := config.NewConfig()
	conf.PrintValues()

	appCtx, appCtxCancel := context.WithCancel(context.Background())
	defer appCtxCancel()
	group, grpCtx := errgroup.WithContext(appCtx)

	logger, _ := zap.NewProduction()

	packer := services.NewPackService()
	packController := controllers.NewPackController(packer)

	httpSrv := api.NewServer(
		grpCtx,
		conf,
		api.WithPackerService(packer),
		api.WithLogger(logger),
	)
	group.Go(func() error {
		if err := httpSrv.Run(
			packController,
		); err != nil {
			logger.Error("error running http server: ", zap.Error(err))
			return err
		}

		httpSrv.Wait()
		return nil
	})

	if err := group.Wait(); err != nil {
		logger.Error("process terminated: ", zap.Error(err))
		return
	}
}
