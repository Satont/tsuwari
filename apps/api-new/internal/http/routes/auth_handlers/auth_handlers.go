package auth_handlers

import (
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/satont/tsuwari/apps/api-new/internal/grpc_clients"
	"github.com/satont/tsuwari/apps/api-new/internal/http/fiber"
	"github.com/satont/tsuwari/apps/api-new/internal/http/middlewares"
	config "github.com/satont/tsuwari/libs/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AuthHandlers struct {
	config         *config.Config
	middlewares    *middlewares.Middlewares
	grpcClients    *grpc_clients.GrpcClients
	logger         *zap.SugaredLogger
	sessionStorage *session.Store
	cacheStorage   *fiber.RedisCacheStorage
}

type Opts struct {
	fx.In

	Config         *config.Config
	Middlewares    *middlewares.Middlewares
	GrpcClients    *grpc_clients.GrpcClients
	Logger         *zap.SugaredLogger
	SessionStorage *session.Store
	CacheStorage   *fiber.RedisCacheStorage
}

func NewAuthHandlers(opts Opts) *AuthHandlers {
	return &AuthHandlers{
		config:         opts.Config,
		middlewares:    opts.Middlewares,
		grpcClients:    opts.GrpcClients,
		logger:         opts.Logger,
		sessionStorage: opts.SessionStorage,
		cacheStorage:   opts.CacheStorage,
	}
}
