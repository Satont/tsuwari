package client

import (
	"context"
	"github.com/dnsge/twitch-eventsub-framework"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/eventsub/internal/creds"
	"github.com/satont/tsuwari/apps/eventsub/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/twitch"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
)

type SubClient struct {
	*eventsub_framework.SubClient

	services    *types.Services
	callbackUrl string
}

func NewClient(
	ctx context.Context,
	services *types.Services,
	callBackUrl string,
) (*SubClient, error) {
	eventSubCreds := creds.NewCreds(ctx, services.Config, services.Grpc.Tokens)
	client := eventsub_framework.NewSubClient(eventSubCreds)

	subClient := &SubClient{
		SubClient:   client,
		services:    services,
		callbackUrl: callBackUrl,
	}

	var channels []model.Channels
	err := services.Gorm.Where(`"isEnabled" = ?`, true).Find(&channels).Error
	if err != nil {
		return nil, err
	}

	for _, channel := range channels {
		err = subClient.SubscribeToNeededEvents(ctx, channel.ID)
		if err != nil {
			return nil, err
		}
	}

	return subClient, nil
}

func (c *SubClient) SubscribeToNeededEvents(ctx context.Context, userId string) error {
	channelCondition := map[string]string{
		"broadcaster_user_id": userId,
	}
	userCondition := map[string]string{
		"user_id": userId,
	}

	neededSubs := map[string]map[string]string{
		"channel.update": channelCondition,
		"stream.online":  channelCondition,
		"stream.offline": channelCondition,
		"user.update":    userCondition,
		"channel.follow": {
			"broadcaster_user_id": userId,
			"moderator_user_id":   userId,
		},
		"channel.moderator.add":                                  channelCondition,
		"channel.moderator.remove":                               channelCondition,
		"channel.channel_points_custom_reward_redemption.add":    channelCondition,
		"channel.channel_points_custom_reward_redemption.update": channelCondition,
	}

	twitchClient, err := twitch.NewAppClient(*c.services.Config, c.services.Grpc.Tokens)
	if err != nil {
		return err
	}

	subsReq, err := twitchClient.GetEventSubSubscriptions(&helix.EventSubSubscriptionsParams{
		UserID: userId,
	})

	wg := &sync.WaitGroup{}
	wg.Add(len(neededSubs))

	var ops uint64

	for key, value := range neededSubs {
		go func(key string, value map[string]string) {
			defer wg.Done()

			existedSub, ok := lo.Find(subsReq.Data.EventSubSubscriptions, func(item helix.EventSubSubscription) bool {
				return item.Type == key && (item.Condition.BroadcasterUserID == value["broadcaster_user_id"] || item.Condition.UserID == value["user_id"])
			})

			if ok && existedSub.Status == "enabled" && existedSub.Transport.Callback == c.callbackUrl {
				return
			}

			if ok {
				err = c.Unsubscribe(ctx, existedSub.ID)
				if err != nil {
					zap.S().Errorw("Failed to unsubscribe", "user_id", userId, "type", key, "error", err)
					return
				}
			}

			request := &eventsub_framework.SubRequest{
				Type:      key,
				Condition: value,
				Callback:  c.callbackUrl,
				Secret:    c.services.Config.TwitchClientSecret,
			}
			if _, err := c.Subscribe(ctx, request); err != nil {
				zap.S().Error(err, key, userId)
				return
			}

			atomic.AddUint64(&ops, 1)
		}(key, value)
	}
	wg.Wait()

	zap.S().Infow("Subcribed to needed events", "user_id", userId, "madeRequests", ops)

	return nil
}
