package gorm

import (
	"context"
	"time"

	cfg "github.com/satont/twir/libs/config"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Config cfg.Config
	LC     fx.Lifecycle
}

func New(opts Opts) (*gorm.DB, error) {
	db, err := gorm.Open(
		postgres.Open(opts.Config.DatabaseUrl),
	)
	if err != nil {
		return nil, err
	}
	d, _ := db.DB()
	d.SetMaxIdleConns(1)
	d.SetMaxOpenConns(10)
	d.SetConnMaxLifetime(time.Hour)

	opts.LC.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				return d.Close()
			},
		},
	)

	return db, nil
}
