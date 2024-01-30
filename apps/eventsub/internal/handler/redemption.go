package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
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
)

func (c *Handler) handleChannelPointsRewardRedemptionAdd(
	h *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelPointsRewardRedemptionAdd,
) {
	c.logger.Info(
		"channel points reward redemption add",
		slog.String("reward", event.Reward.Title),
		slog.String("userName", event.UserLogin),
		slog.String("userId", event.UserID),
		slog.String("channelName", event.BroadcasterUserLogin),
		slog.String("channelId", event.BroadcasterUserID),
	)

	err := c.gorm.Create(
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
	).Error
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}

	// fire event to events microsevice
	_, err = c.eventsGrpc.RedemptionCreated(
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
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}

	// update user spend points
	go func() {
		e := c.countUserChannelPoints(event.UserID, event.BroadcasterUserID, event.Reward.Cost)
		if e != nil {
			c.logger.Error(e.Error(), slog.Any("err", e))
		}
	}()

	// youtube song requests

	go func() {
		e := c.handleYoutubeSongRequests(event)
		if e != nil {
			c.logger.Error(e.Error(), slog.Any("e", err))
		}
	}()

	go func() {
		e := c.handleAlerts(event)
		if e != nil {
			c.logger.Error(e.Error(), slog.Any("e", err))
		}
	}()

	go func() {
		e := c.handleRewardsSevenTvEmote(event)
		if e != nil {
			c.logger.Error(e.Error(), slog.Any("err", e))
		}
	}()

}

func (c *Handler) handleChannelPointsRewardRedemptionUpdate(
	h *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelPointsRewardRedemptionUpdate,
) {
	if event.Status != "CANCELED" {
		return
	}

	userStats := &model.UsersStats{}
	err := c.gorm.Where(`"userId" = ?`, event.UserID).Find(userStats).Error
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
		return
	}
	if userStats.ID == "" {
		return
	}
	userStats.UsedChannelPoints -= int64(event.Reward.Cost)
	err = c.gorm.Save(userStats).Error
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
		return
	}
}

func (c *Handler) countUserChannelPoints(userId, channelId string, count int) error {
	user := &model.Users{}
	err := c.gorm.
		Where("id = ?", userId).
		Preload("Stats", `"channelId" = ?`, channelId).
		First(user).Error
	if err != nil {
		return err
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

		err = c.gorm.Error
		if err != nil {
			return err
		}
	}

	if user.Stats != nil {
		user.Stats.UsedChannelPoints += int64(count)
		err = c.gorm.Save(user.Stats).Error
		if err != nil {
			return err
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
		err = c.gorm.Create(user.Stats).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Handler) handleYoutubeSongRequests(
	event *eventsub_bindings.EventChannelPointsRewardRedemptionAdd,
) error {
	if event.UserInput == "" {
		return nil
	}

	settings := &modules.YouTubeSettings{}
	entity := &model.ChannelModulesSettings{}
	err := c.gorm.
		Where(`"channelId" = ? AND "type" = ?`, event.BroadcasterUserID, "youtube_song_requests").
		Find(entity).
		Error
	if err != nil {
		return err
	}
	if entity.ID == "" {
		return nil
	}

	err = json.Unmarshal(entity.Settings, settings)
	if err != nil {
		return err
	}

	if !*settings.Enabled || event.Reward.ID != settings.ChannelPointsRewardId {
		return nil
	}

	command := &model.ChannelsCommands{}
	err = c.gorm.
		Where(`"channelId" = ? AND "defaultName" = ?`, event.BroadcasterUserID, "sr").
		Find(command).Error
	if err != nil {
		return err
	}
	if command.ID == "" {
		c.logger.Warn("no command sr", slog.String("channelId", event.BroadcasterUserID))
		return nil
	}

	res, err := c.parserGrpc.ProcessCommand(
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
		return err
	}

	if len(res.GetResponses()) == 0 {
		return nil
	}

	for _, response := range res.GetResponses() {
		c.botsGrpc.SendMessage(
			context.Background(),
			&bots.SendMessageRequest{
				ChannelId:   event.BroadcasterUserID,
				ChannelName: &event.BroadcasterUserLogin,
				Message:     fmt.Sprintf("@%s %s", event.UserLogin, response),
				IsAnnounce:  lo.ToPtr(false),
			},
		)
	}

	return nil
}

func (c *Handler) handleAlerts(
	event *eventsub_bindings.EventChannelPointsRewardRedemptionAdd,
) error {
	alert := model.ChannelAlert{}

	if err := c.gorm.Where(
		"channel_id = ? AND reward_ids && ?",
		event.BroadcasterUserID,
		pq.StringArray{event.Reward.ID},
	).Find(&alert).Error; err != nil {
		return err
	}

	if alert.ID == "" {
		return nil
	}

	_, err := c.websocketsGrpc.TriggerAlert(
		context.TODO(),
		&websockets.TriggerAlertRequest{
			ChannelId: event.BroadcasterUserID,
			AlertId:   alert.ID,
		},
	)

	return err
}
