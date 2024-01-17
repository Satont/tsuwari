package activity

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/timers/internal/repositories/channels"
	"github.com/satont/twir/apps/timers/internal/repositories/streams"
	"github.com/satont/twir/apps/timers/internal/repositories/timers"
	config "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/libs/grpc/bots"
	"github.com/twirapp/twir/libs/grpc/parser"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	TimersRepository   timers.Repository
	ChannelsRepository channels.Repository
	StreamsRepository  streams.Repository
	Cfg                config.Config
	ParserGrpc         parser.ParserClient
	BotsGrpc           bots.BotsClient
	Redis              *redis.Client
}

func New(opts Opts) *Activity {
	return &Activity{
		timersRepository:   opts.TimersRepository,
		channelsRepository: opts.ChannelsRepository,
		streamsRepository:  opts.StreamsRepository,
		cfg:                opts.Cfg,
		parserGrpc:         opts.ParserGrpc,
		botsGrpc:           opts.BotsGrpc,
		redis:              opts.Redis,
	}
}

type Activity struct {
	timersRepository   timers.Repository
	channelsRepository channels.Repository
	streamsRepository  streams.Repository
	cfg                config.Config
	parserGrpc         parser.ParserClient
	botsGrpc           bots.BotsClient
	redis              *redis.Client
}

func (c *Activity) SendMessage(ctx context.Context, timerId string, _ int) (
	int,
	error,
) {
	timer, err := c.timersRepository.GetById(timerId)
	if err != nil {
		return 0, err
	}

	channel, err := c.channelsRepository.GetById(timer.ChannelID)
	if err != nil {
		return 0, err
	}

	if !channel.Enabled {
		return 0, nil
	}

	if !channel.IsBotMod {
		return 0, nil
	}

	currentResponse, err := c.redis.Get(ctx, "timers:current_response:"+timerId).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		return currentResponse, err
	}

	stream, err := c.streamsRepository.GetByChannelId(channel.ID)
	if err != nil && errors.Is(err, streams.NotFound) && c.cfg.AppEnv != "development" {
		return currentResponse, nil
	} else if err != nil && c.cfg.AppEnv != "development" {
		return currentResponse, err
	}

	if timer.MessageInterval != 0 &&
		timer.LastTriggerMessageNumber-stream.ParsedMessages+timer.MessageInterval > 0 {
		return currentResponse, nil
	}

	var response timers.TimerResponse
	for index, r := range timer.Responses {
		if index == currentResponse {
			response = r
			break
		}
	}

	if response.ID == "" || response.Text == "" {
		return currentResponse, nil
	}

	err = c.sendMessage(ctx, stream, channel.ID, response.Text, response.IsAnnounce)
	if err != nil {
		return currentResponse, err
	}

	nextIndex := currentResponse + 1

	if nextIndex >= len(timer.Responses) {
		nextIndex = 0
	}

	err = c.timersRepository.UpdateTriggerMessageNumber(timerId, stream.ParsedMessages)
	if err != nil {
		return nextIndex, err
	}

	err = c.redis.Set(ctx, "timers:current_response:"+timerId, nextIndex, 24*time.Hour).Err()
	if err != nil {
		return nextIndex, err
	}

	return nextIndex, nil
}

func (c *Activity) sendMessage(
	ctx context.Context,
	stream streams.Stream,
	channelId, text string,
	isAnnounce bool,
) error {
	parseReq, err := c.parserGrpc.ParseTextResponse(
		ctx,
		&parser.ParseTextRequestData{
			Sender: &parser.Sender{
				Id:          "",
				Name:        "bot",
				DisplayName: "Bot",
				Badges:      []string{"BROADCASTER"},
			},
			Channel: &parser.Channel{Id: stream.UserID, Name: stream.UserLogin},
			Message: &parser.Message{Text: text},
		},
	)
	if err != nil {
		return err
	}

	for _, message := range parseReq.Responses {
		_, err = c.botsGrpc.SendMessage(
			ctx,
			&bots.SendMessageRequest{
				ChannelId:      channelId,
				ChannelName:    &stream.UserLogin,
				Message:        message,
				IsAnnounce:     &isAnnounce,
				SkipRateLimits: true,
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}
