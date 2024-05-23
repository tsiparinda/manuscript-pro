package repo

import (
	"brucheion/models"
	"context"
	"database/sql"
	"sync"

	"github.com/vedicsociety/platform/config"
	"github.com/vedicsociety/platform/logging"
	"github.com/vedicsociety/platform/services"
)

func RegisterSqlRepositoryService() {

	var db *sql.DB
	var commands *SqlCommands
	// var needInit bool

	loadOnce := sync.Once{}

	//resetOnce := sync.Once{}
	services.AddScoped(func(ctx context.Context, config config.Configuration,
		logger logging.Logger) models.Repository {
		loadOnce.Do(func() {

			db, commands = openDB(config, logger)

		})
		repo := &SqlRepository{
			Configuration: config,
			Logger:        logger,
			Commands:      *commands,
			DB:            db,
			Context:       ctx,
		}
		return repo
	})
}
