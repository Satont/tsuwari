package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/api/internal/files"
	"github.com/satont/twir/apps/api/internal/handlers"
	"github.com/satont/twir/apps/api/internal/impl_protected"
	"github.com/satont/twir/apps/api/internal/impl_unprotected"
	"github.com/satont/twir/apps/api/internal/interceptors"
	"github.com/satont/twir/apps/api/internal/sessions"
	"github.com/satont/twir/apps/api/internal/twirp_handlers"
	"github.com/satont/twir/apps/api/internal/webhooks"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	internalSentry "github.com/satont/twir/libs/sentry"
	"github.com/twirapp/twir/libs/grpc/bots"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/discord"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/eventsub"
	"github.com/twirapp/twir/libs/grpc/integrations"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/scheduler"
	"github.com/twirapp/twir/libs/grpc/timers"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var App = fx.Options(
	fx.Provide(
		func() cfg.Config {
			config, err := cfg.New()
			if err != nil {
				panic(err)
			}
			return *config
		},
		internalSentry.NewFx(
			internalSentry.NewFxOpts{
				Service: "api",
			},
		),
		logger.NewFx(logger.Opts{Service: "api"}),
		func(c cfg.Config) tokens.TokensClient {
			return clients.NewTokens(c.AppEnv)
		},
		func(c cfg.Config) bots.BotsClient {
			return clients.NewBots(c.AppEnv)
		},
		func(c cfg.Config) integrations.IntegrationsClient {
			return clients.NewIntegrations(c.AppEnv)
		},
		func(c cfg.Config) parser.ParserClient {
			return clients.NewParser(c.AppEnv)
		},
		func(c cfg.Config) events.EventsClient {
			return clients.NewEvents(c.AppEnv)
		},
		func(c cfg.Config) websockets.WebsocketClient {
			return clients.NewWebsocket(c.AppEnv)
		},
		func(c cfg.Config) scheduler.SchedulerClient {
			return clients.NewScheduler(c.AppEnv)
		},
		func(c cfg.Config) timers.TimersClient {
			return clients.NewTimers(c.AppEnv)
		},
		func(c cfg.Config) eventsub.EventSubClient {
			return clients.NewEventSub(c.AppEnv)
		},
		func(c cfg.Config) discord.DiscordClient {
			return clients.NewDiscord(c.AppEnv)
		},
		func(config cfg.Config, lc fx.Lifecycle) (*redis.Client, error) {
			redisOpts, err := redis.ParseURL(config.RedisUrl)
			if err != nil {
				return nil, err
			}
			client := redis.NewClient(redisOpts)
			lc.Append(
				fx.Hook{
					OnStop: func(ctx context.Context) error {
						return client.Close()
					},
				},
			)

			return client, nil
		},
		func(r *redis.Client) *scs.SessionManager {
			return sessions.New(r)
		},
		func(config cfg.Config, lc fx.Lifecycle) (*gorm.DB, error) {
			db, err := gorm.Open(
				postgres.Open(config.DatabaseUrl),
			)
			if err != nil {
				return nil, err
			}
			d, _ := db.DB()
			d.SetMaxOpenConns(5)
			d.SetConnMaxIdleTime(1 * time.Minute)

			lc.Append(
				fx.Hook{
					OnStop: func(_ context.Context) error {
						return d.Close()
					},
				},
			)

			return db, nil
		},
		interceptors.New,
		impl_protected.New,
		impl_unprotected.New,
		handlers.AsHandler(twirp_handlers.NewProtected),
		handlers.AsHandler(twirp_handlers.NewUnProtected),
		handlers.AsHandler(webhooks.NewDonateStream),
		handlers.AsHandler(webhooks.NewDonatello),
		handlers.AsHandler(files.NewFiles),
		fx.Annotate(
			func(handlers []handlers.IHandler) *http.ServeMux {
				mux := http.NewServeMux()
				for _, route := range handlers {
					mux.Handle(route.Pattern(), route.Handler())
				}
				return mux
			},
			fx.ParamTags(`group:"handlers"`),
		),
	),
	fx.Invoke(
		func(
			mux *http.ServeMux,
			sessionManager *scs.SessionManager,
			l logger.Logger,
			lc fx.Lifecycle,
		) error {
			server := &http.Server{
				Addr:    "0.0.0.0:3002",
				Handler: sessionManager.LoadAndSave(mux),
			}

			lc.Append(
				fx.Hook{
					OnStart: func(_ context.Context) error {
						go func() {
							l.Info("Started", slog.String("port", "3002"))
							err := server.ListenAndServe()
							if err != nil && !errors.Is(err, http.ErrServerClosed) {
								panic(err)
							}
						}()

						return nil
					},
					OnStop: func(ctx context.Context) error {
						return server.Shutdown(ctx)
					},
				},
			)

			return nil
		},
	),
)
