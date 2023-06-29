package users

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/di"
	"github.com/satont/twir/apps/api/internal/interfaces"
	"github.com/satont/twir/apps/api/internal/types"
	config "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/twitch"
	"github.com/satont/twir/libs/types/types/api/modules"
	"go.uber.org/zap"
)

type UserSettings struct {
	Rate  int    `json:"rate"`
	Voice string `json:"voice"`
	Pitch int    `json:"pitch"`

	UserLogin  string `json:"userLogin"`
	UserName   string `json:"userName"`
	UserAvatar string `json:"userAvatar"`
	UserID     string `json:"userId"`
}

func handleGet(channelId string, services types.Services) ([]*UserSettings, error) {
	var settings []model.ChannelModulesSettings
	err := services.DB.
		Where(`"channelId" = ? AND "type" = ? AND "userId" IS NOT NULL`, channelId, "tts").
		Find(&settings).
		Error
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Internal error")
	}

	var usersSettings []*UserSettings

	for _, setting := range settings {
		var ttsSettings modules.TTSSettings
		err = json.Unmarshal(setting.Settings, &ttsSettings)
		if err != nil {
			zap.S().Error(err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Internal error")
		}

		usersSettings = append(
			usersSettings, &UserSettings{
				Rate:   ttsSettings.Rate,
				Voice:  ttsSettings.Voice,
				Pitch:  ttsSettings.Pitch,
				UserID: setting.UserId.String,
			},
		)
	}

	cfg := do.MustInvoke[config.Config](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)
	twitchClient, err := twitch.NewAppClient(cfg, tokensGrpc)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Internal error")
	}

	chunks := lo.Chunk(usersSettings, 100)
	wg := &sync.WaitGroup{}
	wg.Add(len(chunks))

	for _, chunk := range chunks {
		go func(c []*UserSettings) {
			defer wg.Done()

			users, err := twitchClient.GetUsers(
				&helix.UsersParams{
					IDs: lo.Map(
						c, func(item *UserSettings, _ int) string {
							return item.UserID
						},
					),
				},
			)

			if err != nil || users.ErrorMessage != "" {
				zap.S().Error(err, users.ErrorMessage)
				return
			}

			for _, user := range users.Data.Users {
				settings, ok := lo.Find(
					usersSettings, func(s *UserSettings) bool {
						return s.UserID == user.ID
					},
				)
				if !ok {
					continue
				}

				settings.UserAvatar = user.ProfileImageURL
				settings.UserLogin = user.Login
				settings.UserName = user.DisplayName
			}
		}(chunk)
	}

	wg.Wait()

	return usersSettings, nil
}

func handleDelete(channelId string, dto *deleteDto, services types.Services) error {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	user := &model.ChannelModulesSettings{}
	err := services.DB.
		Where(`"userId" IN ? AND "channelId" = ? AND type = ?`, dto.UsersIDS, channelId, "tts").
		Delete(user).
		Error
	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return nil
}
