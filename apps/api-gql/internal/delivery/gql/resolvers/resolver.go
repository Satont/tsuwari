package resolvers

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	config "github.com/satont/twir/libs/config"
	deprecatedgormmodel "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	twir_stats "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/twir-stats"
	admin_actions "github.com/twirapp/twir/apps/api-gql/internal/services/admin-actions"
	"github.com/twirapp/twir/apps/api-gql/internal/services/alerts"
	audit_logs "github.com/twirapp/twir/apps/api-gql/internal/services/audit-logs"
	"github.com/twirapp/twir/apps/api-gql/internal/services/badges"
	badges_users "github.com/twirapp/twir/apps/api-gql/internal/services/badges-users"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands_responses"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands_with_groups_and_responses"
	dashboard_widget_events "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard-widget-events"
	"github.com/twirapp/twir/apps/api-gql/internal/services/greetings"
	"github.com/twirapp/twir/apps/api-gql/internal/services/keywords"
	"github.com/twirapp/twir/apps/api-gql/internal/services/roles"
	"github.com/twirapp/twir/apps/api-gql/internal/services/roles_users"
	"github.com/twirapp/twir/apps/api-gql/internal/services/roles_with_roles_users"
	"github.com/twirapp/twir/apps/api-gql/internal/services/timers"
	ttsvoices "github.com/twirapp/twir/apps/api-gql/internal/services/tts_voices"
	twir_users "github.com/twirapp/twir/apps/api-gql/internal/services/twir-users"
	twitchservice "github.com/twirapp/twir/apps/api-gql/internal/services/twitch"
	"github.com/twirapp/twir/apps/api-gql/internal/services/users"
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

type Deps struct {
	fx.In

	Sessions             *auth.Auth
	Gorm                 *gorm.DB
	Config               config.Config
	TokensGrpc           tokens.TokensClient
	CachedTwitchClient   *twitchcahe.CachedTwitchClient
	CachedCommandsClient *generic_cacher.GenericCacher[[]deprecatedgormmodel.ChannelsCommands]
	Minio                *minio.Client
	TwirBus              *bus_core.Bus
	Logger               logger.Logger
	Redis                *redis.Client
	WsRouter             wsrouter.WsRouter
	TwirStats            *twir_stats.TwirStats

	DashboardWidgetEventsService          *dashboard_widget_events.Service
	VariablesService                      *variables.Service
	TimersService                         *timers.Service
	KeywordsService                       *keywords.Service
	AuditLogsService                      *audit_logs.Service
	AdminActionsService                   *admin_actions.Service
	BadgesService                         *badges.Service
	BadgesUsersService                    *badges_users.Service
	UsersService                          *users.Service
	TwirUsersService                      *twir_users.Service
	AlertsService                         *alerts.Service
	CommandsService                       *commands.Service
	CommandsWithGroupsAndResponsesService *commands_with_groups_and_responses.Service
	CommandsResponsesService              *commands_responses.Service
	GreetingsService                      *greetings.Service
	RolesService                          *roles.Service
	RolesUsersService                     *roles_users.Service
	RolesWithUsersService                 *roles_with_roles_users.Service
	TwitchService                         *twitchservice.Service
	TTSVoicesService                      *ttsvoices.Service
}

type Resolver struct {
	deps Deps
}

func New(deps Deps) (*Resolver, error) {
	return &Resolver{
		deps: deps,
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
