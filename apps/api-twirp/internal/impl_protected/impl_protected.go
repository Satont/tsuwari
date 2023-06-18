package impl_protected

import (
	"github.com/redis/go-redis/v9"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected/bot"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected/commands"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected/community"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected/deps"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected/events"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected/greetings"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected/integrations"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected/keywords"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected/modules"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected/rewards"
	"gorm.io/gorm"
)

type Api struct {
	*integrations.Integrations
	*keywords.Keywords
	*modules.Modules
	*bot.Bot
	*commands.Commands
	*community.Community
	*events.Events
	*greetings.Greetings
	*rewards.Rewards
}

type Opts struct {
	Redis *redis.Client
	DB    *gorm.DB
}

func NewApi(opts Opts) *Api {
	d := &deps.Deps{
		Redis: opts.Redis,
		Db:    opts.DB,
	}

	return &Api{
		Integrations: &integrations.Integrations{Deps: d},
		Keywords:     &keywords.Keywords{Deps: d},
		Modules:      &modules.Modules{Deps: d},
		Bot:          &bot.Bot{Deps: d},
		Commands:     &commands.Commands{Deps: d},
		Community:    &community.Community{Deps: d},
		Events:       &events.Events{Deps: d},
		Greetings:    &greetings.Greetings{Deps: d},
		Rewards:      &rewards.Rewards{Deps: d},
	}
}
