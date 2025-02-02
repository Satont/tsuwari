package bus_listener

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/emotes-cacher/internal/emotes"
	"github.com/satont/twir/libs/logger"
	buscore "github.com/twirapp/twir/libs/bus-core"
	emotes_cacher "github.com/twirapp/twir/libs/bus-core/emotes-cacher"
	"go.uber.org/fx"
)

type BusListener struct {
	redis  *redis.Client
	logger logger.Logger
	bus    *buscore.Bus
}

type Opts struct {
	fx.In
	Lc fx.Lifecycle

	Redis  *redis.Client
	Logger logger.Logger
	Bus    *buscore.Bus
}

func New(opts Opts) {
	impl := &BusListener{
		redis:  opts.Redis,
		logger: opts.Logger,
		bus:    opts.Bus,
	}

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := impl.bus.EmotesCacher.CacheGlobalEmotes.SubscribeGroup(
					"emotes-cacher",
					impl.cacheGlobalEmotes,
				); err != nil {
					return err
				}
				if err := impl.bus.EmotesCacher.CacheChannelEmotes.SubscribeGroup(
					"emotes-cacher",
					impl.cacheChannelEmotes,
				); err != nil {
					return err
				}
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return nil
			},
		},
	)
}

func (c *BusListener) cacheChannelEmotes(
	_ context.Context,
	req emotes_cacher.EmotesCacheRequest,
) struct{} {
	if req.ChannelID == "" {
		return struct{}{}
	}

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	resultEmotes := make([]string, 0, 300)

	reqFuncs := []func(c string) ([]string, error){
		emotes.GetChannelSevenTvEmotes,
		emotes.GetChannelBttvEmotes,
		emotes.GetChannelFfzEmotes,
	}

	for _, f := range reqFuncs {
		wg.Add(1)
		f := f
		go func() {
			defer wg.Done()
			res, err := f(req.ChannelID)
			if err != nil {
				c.logger.Error("cannot get emotes", slog.Any("err", err))
				return
			}

			mu.Lock()
			resultEmotes = append(resultEmotes, res...)
			mu.Unlock()
		}()
	}

	wg.Wait()

	c.redis.Pipelined(
		context.Background(), func(pipe redis.Pipeliner) error {
			for _, emote := range resultEmotes {
				if emote == "" {
					continue
				}

				pipe.Set(
					context.Background(),
					fmt.Sprintf("emotes:channel:%s:%s", req.ChannelID, emote),
					emote,
					10*time.Minute,
				)
			}

			return nil
		},
	)

	return struct{}{}
}

func (c *BusListener) cacheGlobalEmotes(_ context.Context, _ struct{}) struct{} {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	resultEmotes := make([]string, 300)

	wg.Add(3)

	go func() {
		defer wg.Done()
		em, err := emotes.GetGlobalSevenTvEmotes()
		if err != nil || em == nil || len(em) == 0 {
			return
		}

		mu.Lock()
		defer mu.Unlock()
		resultEmotes = append(resultEmotes, em...)
	}()

	go func() {
		defer wg.Done()
		em, err := emotes.GetGlobalFfzEmotes()
		if err != nil || em == nil || len(em) == 0 {
			return
		}

		mu.Lock()
		defer mu.Unlock()
		resultEmotes = append(resultEmotes, em...)
	}()

	go func() {
		defer wg.Done()
		em, err := emotes.GetGlobalBttvEmotes()
		if err != nil || em == nil || len(em) == 0 {
			return
		}

		mu.Lock()
		defer mu.Unlock()
		resultEmotes = append(resultEmotes, em...)
	}()

	wg.Wait()

	c.redis.Pipelined(
		context.Background(), func(pipe redis.Pipeliner) error {
			for _, emote := range resultEmotes {
				if emote == "" {
					continue
				}

				pipe.Set(
					context.Background(),
					fmt.Sprintf("emotes:global:%s", emote),
					emote,
					10*time.Minute,
				)
			}

			return nil
		},
	)

	return struct{}{}
}
