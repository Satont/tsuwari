package app

import (
	eventsActivity "github.com/satont/twir/apps/events/internal/activities/events"
	"github.com/satont/twir/apps/events/internal/chat_alerts"
	"github.com/satont/twir/apps/events/internal/grpc_impl"
	"github.com/satont/twir/apps/events/internal/hydrator"
	"github.com/satont/twir/apps/events/internal/workers"
	"github.com/satont/twir/apps/events/internal/workflows"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/grpc/websockets"
	greetingsrepository "github.com/twirapp/twir/libs/repositories/greetings"
	greetingsrepositorypgx "github.com/twirapp/twir/libs/repositories/greetings/pgx"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"
)

var App = fx.Module(
	"events",
	baseapp.CreateBaseApp(baseapp.Opts{AppName: "events"}),
	fx.Provide(
		fx.Annotate(
			greetingsrepositorypgx.NewFx,
			fx.As(new(greetingsrepository.Repository)),
		),
		func(config cfg.Config) tokens.TokensClient {
			return clients.NewTokens(config.AppEnv)
		},
		func(config cfg.Config) websockets.WebsocketClient {
			return clients.NewWebsocket(config.AppEnv)
		},
		hydrator.New,
		eventsActivity.New,
		workflows.NewEventsWorkflow,
		chat_alerts.New,
	),
	fx.Invoke(
		uptrace.NewFx("events"),
		workers.NewEventsWorker,
		grpc_impl.New,
		func(l logger.Logger) {
			l.Info("Events service started")
		},
	),
)
