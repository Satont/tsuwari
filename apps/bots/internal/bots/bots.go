package bots

import (
	"sync"

	cfg "github.com/satont/tsuwari/libs/config"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/satont/tsuwari/libs/twitch"

	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/bots/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type NewBotsOpts struct {
	Twitch *twitch.Twitch
	DB     *gorm.DB
	Logger *zap.Logger
	Cfg    *cfg.Config
	Nats   *nats.Conn
}

type BotsService struct {
	Instances map[string]*types.BotClient
}

func NewBotsService(opts *NewBotsOpts) *BotsService {
	service := BotsService{
		Instances: make(map[string]*types.BotClient),
	}
	mu := sync.Mutex{}

	bots := []model.Bots{}
	err := opts.DB.
		Preload("Token").
		Preload("Channels").
		Find(&bots).
		Error
	if err != nil {
		panic(err)
	}

	for _, bot := range bots {
		go func(bot model.Bots) {
			instance := newBot(&ClientOpts{
				DB:     opts.DB,
				Cfg:    opts.Cfg,
				Logger: opts.Logger,
				Bot:    &bot,
				Nats:   opts.Nats,
			})

			channels := []model.Channels{}
			opts.DB.Where(`"botId" = ?`, bot.ID).Select("ID", "BotID").Find(&channels)

			if len(channels) > 0 {
				ids := lo.Map(channels, func(i model.Channels, _ int) string {
					return i.ID
				})

				opts.DB.Model(&model.ChannelsGreetings{}).
					Where(`"channelId" IN ?`, ids).
					Update("processed", false)
			}

			mu.Lock()
			service.Instances[bot.ID] = instance
			mu.Unlock()
		}(bot)
	}

	return &service
}
