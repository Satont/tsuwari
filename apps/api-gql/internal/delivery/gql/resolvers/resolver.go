package resolvers

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/minio/minio-go/v7"
	"github.com/nicklaw5/helix/v2"
	"github.com/redis/go-redis/v9"
	config "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	auditlogs "github.com/satont/twir/libs/pubsub/audit-logs"
	"github.com/satont/twir/libs/twitch"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	twir_stats "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/twir-stats"
	dashboard_widget_events "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard-widget-events"
	"github.com/twirapp/twir/apps/api-gql/internal/services/keywords"
	"github.com/twirapp/twir/apps/api-gql/internal/services/timers"
	"github.com/twirapp/twir/apps/api-gql/internal/services/variables"
	"github.com/twirapp/twir/apps/api-gql/internal/wsrouter"
	bus_core "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	twitchcahe "github.com/twirapp/twir/libs/cache/twitch"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	config               config.Config
	sessions             *auth.Auth
	gorm                 *gorm.DB
	twitchClient         *helix.Client
	cachedTwitchClient   *twitchcahe.CachedTwitchClient
	cachedCommandsClient *generic_cacher.GenericCacher[[]model.ChannelsCommands]
	minioClient          *minio.Client
	twirBus              *bus_core.Bus
	logger               logger.Logger
	redis                *redis.Client
	keywordsCacher       *generic_cacher.GenericCacher[[]model.ChannelsKeywords]
	tokensClient         tokens.TokensClient
	wsRouter             wsrouter.WsRouter
	auditLogsPubSub      auditlogs.PubSub
	twirStats            *twir_stats.TwirStats

	dashboardWidgetEventsService *dashboard_widget_events.DashboardWidgetsEvents
	variablesService             *variables.Service
	timersService                *timers.Service
	keywordsService              *keywords.Service
}

type Opts struct {
	fx.In

	Sessions             *auth.Auth
	Gorm                 *gorm.DB
	Config               config.Config
	TokensGrpc           tokens.TokensClient
	CachedTwitchClient   *twitchcahe.CachedTwitchClient
	CachedCommandsClient *generic_cacher.GenericCacher[[]model.ChannelsCommands]
	Minio                *minio.Client
	TwirBus              *bus_core.Bus
	Logger               logger.Logger
	Redis                *redis.Client
	KeywordsCacher       *generic_cacher.GenericCacher[[]model.ChannelsKeywords]
	WsRouter             wsrouter.WsRouter
	AuditLogsPubSub      auditlogs.PubSub
	TwirStats            *twir_stats.TwirStats

	DashboardWidgetEventsService *dashboard_widget_events.DashboardWidgetsEvents
	VariablesService             *variables.Service
	TimersService                *timers.Service
	KeywordService               *keywords.Service
}

func New(opts Opts) (*Resolver, error) {
	twitchClient, err := twitch.NewAppClient(opts.Config, opts.TokensGrpc)
	if err != nil {
		return nil, err
	}

	return &Resolver{
		config:                       opts.Config,
		sessions:                     opts.Sessions,
		gorm:                         opts.Gorm,
		twitchClient:                 twitchClient,
		cachedTwitchClient:           opts.CachedTwitchClient,
		minioClient:                  opts.Minio,
		twirBus:                      opts.TwirBus,
		logger:                       opts.Logger,
		redis:                        opts.Redis,
		cachedCommandsClient:         opts.CachedCommandsClient,
		keywordsCacher:               opts.KeywordsCacher,
		tokensClient:                 opts.TokensGrpc,
		wsRouter:                     opts.WsRouter,
		auditLogsPubSub:              opts.AuditLogsPubSub,
		twirStats:                    opts.TwirStats,
		dashboardWidgetEventsService: opts.DashboardWidgetEventsService,
		variablesService:             opts.VariablesService,
		timersService:                opts.TimersService,
		keywordsService:              opts.KeywordService,
	}, nil
}

func GetPreloads(ctx context.Context) []string {
	return GetNestedPreloads(
		graphql.GetOperationContext(ctx),
		graphql.CollectFieldsCtx(ctx, nil),
		"",
	)
}

func GetNestedPreloads(
	ctx *graphql.OperationContext,
	fields []graphql.CollectedField,
	prefix string,
) (preloads []string) {
	for _, column := range fields {
		prefixColumn := GetPreloadString(prefix, column.Name)
		preloads = append(preloads, prefixColumn)
		preloads = append(
			preloads,
			GetNestedPreloads(ctx, graphql.CollectFields(ctx, column.Selections, nil), prefixColumn)...,
		)
	}
	return
}

func GetPreloadString(prefix, name string) string {
	if len(prefix) > 0 {
		return prefix + "." + name
	}
	return name
}
