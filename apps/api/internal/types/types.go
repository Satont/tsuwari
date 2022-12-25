package types

import (
	"github.com/satont/tsuwari/libs/grpc/generated/bots"
	"github.com/satont/tsuwari/libs/grpc/generated/eventsub"
	"github.com/satont/tsuwari/libs/grpc/generated/integrations"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"
	"github.com/satont/tsuwari/libs/grpc/generated/scheduler"
	"github.com/satont/tsuwari/libs/grpc/generated/timers"
	"github.com/satont/tsuwari/libs/twitch"

	cfg "github.com/satont/tsuwari/libs/config"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/satont/tsuwari/apps/api/internal/services/redis"
	"gorm.io/gorm"
)

type Services struct {
	DB                  *gorm.DB
	RedisStorage        *redis.RedisStorage
	Validator           *validator.Validate
	ValidatorTranslator ut.Translator
	Twitch              *twitch.Twitch
	Cfg                 *cfg.Config
	TgBotApi            *tgbotapi.BotAPI
	BotsGrpc            bots.BotsClient
	TimersGrpc          timers.TimersClient
	SchedulerGrpc       scheduler.SchedulerClient
	ParserGrpc          parser.ParserClient
	EventSubGrpc        eventsub.EventSubClient
	IntegrationsGrpc    integrations.IntegrationsClient
}

type JSONResult struct{}

type DOCApiBadRequest struct {
	Messages string
}

type DOCApiValidationError struct {
	Messages []string
}

type DOCApiInternalError struct {
	Messages []string
}
