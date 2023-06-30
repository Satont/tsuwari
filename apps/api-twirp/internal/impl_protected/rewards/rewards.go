package rewards

import (
	"context"
	"fmt"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api-twirp/internal/impl_deps"
	"github.com/satont/twir/libs/grpc/generated/api/rewards"
	"github.com/satont/twir/libs/twitch"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Rewards struct {
	*impl_deps.Deps
}

func (c *Rewards) RewardsGet(
	ctx context.Context,
	_ *emptypb.Empty,
) (*rewards.GetResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	twitchClient, err := twitch.NewUserClientWithContext(ctx, dashboardId, *c.Config, c.Grpc.Tokens)
	if err != nil {
		return nil, err
	}

	resp, err := twitchClient.GetCustomRewards(&helix.GetCustomRewardsParams{
		BroadcasterID: dashboardId,
	})
	if err != nil {
		return nil, err
	}
	if resp.ErrorMessage != "" {
		return nil, fmt.Errorf("cannot get channel rewards: %v %s", resp.StatusCode, resp.ErrorMessage)
	}

	return &rewards.GetResponse{
		Rewards: lo.Map(
			resp.Data.ChannelCustomRewards,
			func(item helix.ChannelCustomReward, _ int) *rewards.Reward {
				return &rewards.Reward{
					Id:     item.ID,
					Title:  item.Title,
					Prompt: item.Prompt,
					Cost:   uint64(item.Cost),
					Image: &rewards.Reward_Image{
						Url_1X: item.Image.Url1x,
						Url_2X: item.Image.Url2x,
						Url_4X: item.Image.Url4x,
					},
					BackgroundColor: item.BackgroundColor,
					DefaultImage: &rewards.Reward_Image{
						Url_1X: item.DefaultImage.Url1x,
						Url_2X: item.DefaultImage.Url2x,
						Url_4X: item.DefaultImage.Url4x,
					},
					IsEnabled:           item.IsEnabled,
					IsUserInputRequired: item.IsUserInputRequired,
					MaxPerStreamSetting: &rewards.Reward_MaxPerStreamSettings{
						IsEnabled:    item.MaxPerStreamSetting.IsEnabled,
						MaxPerStream: uint64(item.MaxPerStreamSetting.MaxPerStream),
					},
					MaxPerUserPerStreamSetting: &rewards.Reward_MaxPerUserPerStreamSettings{
						IsEnabled:    item.MaxPerUserPerStreamSetting.IsEnabled,
						MaxPerStream: uint64(item.MaxPerUserPerStreamSetting.MaxPerUserPerStream),
					},
					GlobalCooldownSetting: &rewards.Reward_GlobalCooldownSettings{
						IsEnabled:             item.GlobalCooldownSetting.IsEnabled,
						GlobalCooldownSeconds: uint64(item.GlobalCooldownSetting.GlobalCooldownSeconds),
					},
					IsPaused:                          item.IsPaused,
					IsInStock:                         item.IsInStock,
					ShouldRedemptionsSkipRequestQueue: item.ShouldRedemptionsSkipRequestQueue,
					RedemptionsRedeemedCurrentStream:  uint64(item.RedemptionsRedeemedCurrentStream),
					CooldownExpiresAt:                 item.CooldownExpiresAt,
				}
			},
		),
	}, nil
}
