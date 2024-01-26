package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/lib/pq"
	"github.com/twirapp/twir/libs/grpc/websockets"

	eventsub_bindings "github.com/dnsge/twitch-eventsub-bindings"
	"github.com/google/uuid"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/types/types/api/modules"
	"github.com/twirapp/twir/libs/grpc/bots"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/parser"
	"go.uber.org/zap"
)

func (c *Handler) handleChannelPointsRewardRedemptionAdd(
	h *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelPointsRewardRedemptionAdd,
) {
	zap.S().Infow(
		"channel points reward redemption add",
		"reward", event.Reward.Title,
		"userName", event.UserLogin,
		"userId", event.UserID,
		"channelName", event.BroadcasterUserLogin,
		"channelId", event.BroadcasterUserID,
	)

	c.services.Gorm.Create(
		&model.ChannelsEventsListItem{
			ID:        uuid.New().String(),
			ChannelID: event.BroadcasterUserID,
			UserID:    event.UserID,
			Type:      model.ChannelEventListItemTypeRedemptionCreated,
			CreatedAt: time.Now().UTC(),
			Data: &model.ChannelsEventsListItemData{
				RedemptionInput:           event.UserInput,
				RedemptionTitle:           event.Reward.Title,
				RedemptionUserName:        event.UserLogin,
				RedemptionUserDisplayName: event.UserName,
				RedemptionCost:            strconv.Itoa(event.Reward.Cost),
			},
		},
	)

	// fire event to events microsevice
	c.services.Grpc.Events.RedemptionCreated(
		context.Background(),
		&events.RedemptionCreatedMessage{
			BaseInfo:        &events.BaseInfo{ChannelId: event.BroadcasterUserID},
			UserName:        event.UserLogin,
			UserDisplayName: event.UserName,
			Id:              event.Reward.ID,
			RewardName:      event.Reward.Title,
			RewardCost:      strconv.Itoa(event.Reward.Cost),
			Input:           lo.If(event.UserInput != "", &event.UserInput).Else(nil),
			UserId:          event.UserID,
		},
	)

	// update user spend points
	c.countUserChannelPoints(event.UserID, event.BroadcasterUserID, event.Reward.Cost)

	// youtube song requests
	c.handleYoutubeSongRequests(event)

	c.handleAlerts(event)

	c.handleRewardsSevenTvEmote(event)
}

func (c *Handler) handleChannelPointsRewardRedemptionUpdate(
	h *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelPointsRewardRedemptionUpdate,
) {
	if event.Status != "CANCELED" {
		return
	}

	userStats := &model.UsersStats{}
	err := c.services.Gorm.Where(`"userId" = ?`, event.UserID).Find(userStats).Error
	if err != nil {
		zap.S().Error(err)
		return
	}
	if userStats.ID == "" {
		return
	}
	userStats.UsedChannelPoints -= int64(event.Reward.Cost)
	err = c.services.Gorm.Save(userStats).Error
	if err != nil {
		zap.S().Error(err)
		return
	}
}

func (c *Handler) countUserChannelPoints(userId, channelId string, count int) {
	user := &model.Users{}
	err := c.services.Gorm.
		Where("id = ?", userId).
		Preload("Stats", `"channelId" = ?`, channelId).
		First(user).Error
	if err != nil {
		zap.S().Error(err)
		return
	}

	if user.ID == "" {
		user = &model.Users{
			ID:         "",
			TokenID:    sql.NullString{},
			IsTester:   false,
			IsBotAdmin: false,
			ApiKey:     uuid.New().String(),
			Stats: &model.UsersStats{
				ID:                uuid.New().String(),
				UserID:            userId,
				ChannelID:         channelId,
				Messages:          0,
				Watched:           0,
				UsedChannelPoints: int64(count),
				Emotes:            0,
			},
		}

		err = c.services.Gorm.Create(user).Error
		if err != nil {
			zap.S().Error(err)
			return
		}
	}

	if user.Stats != nil {
		user.Stats.UsedChannelPoints += int64(count)
		err = c.services.Gorm.Save(user.Stats).Error
		if err != nil {
			zap.S().Error(err)
			return
		}
	} else {
		user.Stats = &model.UsersStats{
			ID:                uuid.New().String(),
			UserID:            userId,
			ChannelID:         channelId,
			Messages:          0,
			Watched:           0,
			UsedChannelPoints: int64(count),
			Emotes:            0,
		}
		err = c.services.Gorm.Create(user.Stats).Error
		if err != nil {
			zap.S().Error(err)
			return
		}
	}
}

func (c *Handler) handleYoutubeSongRequests(event *eventsub_bindings.EventChannelPointsRewardRedemptionAdd) {
	if event.UserInput == "" {
		return
	}

	settings := &modules.YouTubeSettings{}
	entity := &model.ChannelModulesSettings{}
	err := c.services.Gorm.
		Where(`"channelId" = ? AND "type" = ?`, event.BroadcasterUserID, "youtube_song_requests").
		Find(entity).
		Error
	if err != nil {
		zap.S().Error(err)
		return
	}
	if entity.ID == "" {
		return
	}

	err = json.Unmarshal(entity.Settings, settings)
	if err != nil {
		zap.S().Error(err)
		return
	}

	if !*settings.Enabled || event.Reward.ID != settings.ChannelPointsRewardId {
		return
	}

	command := &model.ChannelsCommands{}
	err = c.services.Gorm.
		Where(`"channelId" = ? AND "defaultName" = ?`, event.BroadcasterUserID, "sr").
		Find(command).Error
	if err != nil {
		zap.S().Error(err)
		return
	}
	if command.ID == "" {
		zap.S().Warnln("no command sr", "channelId", event.BroadcasterUserID)
		return
	}

	res, err := c.services.Grpc.Parser.ProcessCommand(
		context.Background(), &parser.ProcessCommandRequest{
			Sender: &parser.Sender{
				Id:          event.UserID,
				Name:        event.UserLogin,
				DisplayName: event.UserName,
				Badges:      []string{"VIEWER"},
			},
			Channel: &parser.Channel{
				Id:   event.BroadcasterUserID,
				Name: event.BroadcasterUserName,
			},
			Message: &parser.Message{
				Text:   fmt.Sprintf("!%s %s", command.Name, event.UserInput),
				Id:     event.ID,
				Emotes: nil,
			},
		},
	)

	if err != nil {
		zap.S().Error(err)
		return
	}

	if len(res.GetResponses()) == 0 {
		return
	}

	for _, response := range res.GetResponses() {
		c.services.Grpc.Bots.SendMessage(
			context.Background(),
			&bots.SendMessageRequest{
				ChannelId:   event.BroadcasterUserID,
				ChannelName: &event.BroadcasterUserLogin,
				Message:     fmt.Sprintf("@%s %s", event.UserLogin, response),
				IsAnnounce:  lo.ToPtr(false),
			},
		)
	}

	return
}

func (c *Handler) handleAlerts(event *eventsub_bindings.EventChannelPointsRewardRedemptionAdd) {
	alert := model.ChannelAlert{}

	if err := c.services.Gorm.Where(
		"channel_id = ? AND reward_ids && ?",
		event.BroadcasterUserID,
		pq.StringArray{event.Reward.ID},
	).Find(&alert).Error; err != nil {
		zap.S().Error(err)
		return
	}

	if alert.ID == "" {
		return
	}

	_, err := c.services.Grpc.WebSockets.TriggerAlert(
		context.TODO(),
		&websockets.TriggerAlertRequest{
			ChannelId: event.BroadcasterUserID,
			AlertId:   alert.ID,
		},
	)

	if err != nil {
		zap.S().Error(err)
	}
}
