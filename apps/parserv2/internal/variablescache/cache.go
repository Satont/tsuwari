package variablescache

import (
	"regexp"
	"sync"
	"tsuwari/parser/internal/config/twitch"
	"tsuwari/parser/pkg/helpers"

	"github.com/go-redis/redis/v9"
	"github.com/nicklaw5/helix"
)

type VariablesCacheServices struct {
	Redis  *redis.Client
	Regexp regexp.Regexp
	Twitch *twitch.Twitch
}

type VariablesCacheContext struct {
	ChannelId  string
	SenderId   string
	SenderName string
	Text       string
}

type VariablesCacheService struct {
	Context  VariablesCacheContext
	Services VariablesCacheServices
	Cache    VariablesCache
}

type VariablesCache struct {
	Stream *helix.Stream
}

func New(text string, senderId string, channelId string, senderName *string, redis *redis.Client, r regexp.Regexp, twitch *twitch.Twitch) *VariablesCacheService {
	cache := &VariablesCacheService{
		Context: VariablesCacheContext{
			ChannelId:  channelId,
			SenderId:   senderId,
			SenderName: *senderName,
			Text:       text,
		},
		Services: VariablesCacheServices{
			Redis:  redis,
			Regexp: r,
			Twitch: twitch,
		},
		Cache: VariablesCache{
			Stream: nil,
		},
	}

	cache.fillCache()

	return cache
}

func (c *VariablesCacheService) fillCache() {
	matches := c.Services.Regexp.FindAllStringSubmatch(c.Context.Text, len(c.Context.Text))
	myMap := map[string]interface{}{
		"stream": c.setChannelStream,
	}
	requesting := []string{}
	wg := sync.WaitGroup{}

	c.Services.Twitch.RefreshIfNeeded()

	for _, match := range matches {
		if match[2] == "" {
			continue
		}

		if helpers.Contains(requesting, match[2]) {
			continue
		}

		if val, ok := myMap[match[2]]; ok {
			wg.Add(1)

			go val.(func(wg *sync.WaitGroup))(&wg)
		}
	}

	wg.Wait()
}

func (c *VariablesCacheService) setChannelStream(wg *sync.WaitGroup) {
	defer wg.Done()

	streams, err := c.Services.Twitch.Client.GetStreams(&helix.StreamsParams{
		UserIDs: []string{c.Context.ChannelId},
	})

	if err != nil || streams.Data.Streams == nil {
		return
	}

	c.Cache.Stream = &streams.Data.Streams[0]
}
