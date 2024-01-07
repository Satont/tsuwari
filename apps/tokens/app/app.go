package app

import (
	"github.com/satont/twir/apps/tokens/internal/gorm"
	"github.com/satont/twir/apps/tokens/internal/grpc_impl"
	"github.com/satont/twir/apps/tokens/internal/redis"
	"github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	twirsentry "github.com/satont/twir/libs/sentry"
	"go.uber.org/fx"
)

var App = fx.Module(
	"tokens",
	fx.Provide(
		cfg.NewFx,
		twirsentry.NewFx(twirsentry.NewFxOpts{Service: "tokens"}),
		logger.NewFx(logger.Opts{Service: "tokens"}),
		gorm.New,
		redis.New,
		redis.NewRedisLock,
	),
	fx.Invoke(
		grpc_impl.NewTokens,
		func(l logger.Logger) {
			l.Info("Started")
		},
	),
)
