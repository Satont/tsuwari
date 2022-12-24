package auth

import (
	"context"
	"database/sql"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	"time"

	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/eventsub"
	"github.com/satont/tsuwari/libs/grpc/generated/scheduler"

	"github.com/gofiber/fiber/v2"
	"github.com/satont/go-helix/v2"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func checkUser(
	username, userId string,
	tokens helix.AccessCredentials,
) error {
	db := do.MustInvoke[*gorm.DB](di.Injector)
	logger := do.MustInvoke[interfaces.Logger](di.Injector)

	defaultBot := model.Bots{}
	err := db.Where("type = ?", "DEFAULT").First(&defaultBot).Error
	if err != nil {
		return fiber.NewError(500, "bot not created, cannot create user")
	}

	tokenData := model.Tokens{
		AccessToken:         tokens.AccessToken,
		RefreshToken:        tokens.RefreshToken,
		ExpiresIn:           int32(tokens.ExpiresIn),
		ObtainmentTimestamp: time.Now(),
	}

	user := model.Users{}
	err = db.
		Where(`"users"."id" = ?`, userId).
		Joins("Channel").
		Joins("Token").
		First(&user).Error

	if err != nil && err == gorm.ErrRecordNotFound {
		newToken := tokenData
		newToken.ID = uuid.NewV4().String()

		if err = db.Save(&newToken).Error; err != nil {
			logger.Error(err)
			return err
		}

		user.ID = userId
		user.TokenID = sql.NullString{String: newToken.ID, Valid: true}
		user.ApiKey = uuid.NewV4().String()

		if err = db.Save(&user).Error; err != nil {
			logger.Error(err)
			return err
		}

		channel := createChannelModel(user.ID, defaultBot.ID)
		if err = db.Create(&channel).Error; err != nil {
			logger.Error(err)
			return err
		}
		user.Channel = &channel
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error(err)
		return fiber.NewError(500, "internal error")
	} else {
		if user.Channel == nil {
			channel := createChannelModel(user.ID, defaultBot.ID)
			if err = db.Create(&channel).Error; err != nil {
				logger.Error(err)
				return err
			}
			user.Channel = &channel
		}

		if user.TokenID.Valid {
			tokenData.ID = user.TokenID.String
			if err = db.Select("*").Save(&tokenData).Error; err != nil {
				logger.Error(err)
				return err
			}
		} else {
			tokenData.ID = uuid.NewV4().String()
			if err = db.Save(&tokenData).Error; err != nil {
				logger.Error(err)
				return err
			}
			user.TokenID = sql.NullString{String: tokenData.ID, Valid: true}
			if err := db.Save(&user).Error; err != nil {
				logger.Error(err)
			}
		}
	}

	//if user.Channel.IsEnabled {
	//	services.BotsGrpc.Join(context.Background(), &bots.JoinOrLeaveRequest{
	//		BotId:    user.Channel.BotID,
	//		UserName: username,
	//	})
	//} else {
	//	services.BotsGrpc.Leave(context.Background(), &bots.JoinOrLeaveRequest{
	//		BotId:    user.Channel.BotID,
	//		UserName: username,
	//	})
	//}

	schedulerGrpc := do.MustInvoke[scheduler.SchedulerClient](di.Injector)
	eventSubGrpc := do.MustInvoke[eventsub.EventSubClient](di.Injector)

	schedulerGrpc.CreateDefaultCommands(
		context.Background(),
		&scheduler.CreateDefaultCommandsRequest{
			UserId: userId,
		},
	)
	eventSubGrpc.SubscribeToEvents(
		context.Background(),
		&eventsub.SubscribeToEventsRequest{
			ChannelId: userId,
		},
	)

	return nil
}

func createChannelModel(userId, botId string) model.Channels {
	return model.Channels{
		ID:             userId,
		IsEnabled:      true,
		IsTwitchBanned: false,
		IsBanned:       false,
		BotID:          botId,
	}
}
