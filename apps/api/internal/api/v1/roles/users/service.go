package roles_users

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"

	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/twitch"
	"gorm.io/gorm"
)

type roleUser struct {
	model.ChannelRoleUser
	UserName        string `json:"userName"`
	UserDisplayName string `json:"userDisplayName"`
	UserAvatar      string `json:"userAvatar"`
}

func getUsersService(roleId string, services *types.Services) ([]*roleUser, error) {
	twitchClient, err := twitch.NewAppClient(*services.Config, services.Grpc.Tokens)
	if err != nil {
		services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	role := &model.ChannelRole{}
	if err := services.Gorm.Preload("Users").Where("id = ?", roleId).First(role).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Role not found")
	}

	response := make([]*roleUser, 0, len(role.Users))
	if len(role.Users) == 0 {
		return response, nil
	}

	twitchUsers, err := twitchClient.GetUsers(&helix.UsersParams{
		IDs: lo.Map(role.Users, func(user *model.ChannelRoleUser, _ int) string {
			return user.UserID
		}),
	})

	if err != nil || twitchUsers.ErrorMessage != "" {
		services.Logger.Error(err, twitchUsers.ErrorMessage)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	for _, user := range role.Users {
		twitchUser, ok := lo.Find(twitchUsers.Data.Users, func(twitchUser helix.User) bool {
			return twitchUser.ID == user.UserID
		})

		if !ok {
			continue
		}

		response = append(response, &roleUser{
			ChannelRoleUser: *user,
			UserName:        twitchUser.Login,
			UserDisplayName: twitchUser.DisplayName,
			UserAvatar:      twitchUser.ProfileImageURL,
		})
	}

	return response, nil
}

func updateUsersService(roleId string, userNames []string, services *types.Services) error {
	twitchClient, err := twitch.NewAppClient(*services.Config, services.Grpc.Tokens)
	if err != nil {
		services.Logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	role := &model.ChannelRole{}
	if err := services.Gorm.Preload("Users").Where("id = ?", roleId).First(role).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Role not found")
	}

	if len(userNames) == 0 {
		if err = services.Gorm.Where(`"roleId" = ?`, role.ID).Delete(&model.ChannelRoleUser{}).Error; err != nil {
			services.Logger.Error(err)
			return fiber.NewError(http.StatusInternalServerError, "internal error")
		}

		for _, user := range role.Users {
			services.RedisStorage.DeleteByMethod(
				fmt.Sprintf("fiber:cache:auth:profile:dashboards:%s", user.UserID),
				"GET",
			)
		}

		return nil
	}

	twitchUsers, err := twitchClient.GetUsers(&helix.UsersParams{
		Logins: userNames,
	})

	if err != nil || twitchUsers.ErrorMessage != "" {
		services.Logger.Error(err, twitchUsers.ErrorMessage)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	roleUsers := make([]*model.ChannelRoleUser, 0, len(twitchUsers.Data.Users))
	for _, twitchUser := range twitchUsers.Data.Users {
		roleUsers = append(roleUsers, &model.ChannelRoleUser{
			RoleID: role.ID,
			UserID: twitchUser.ID,
		})
	}

	err = services.Gorm.Transaction(func(tx *gorm.DB) error {
		if err = tx.Where(`"roleId" = ?`, role.ID).Delete(&model.ChannelRoleUser{}).Error; err != nil {
			return err
		}

		if err = tx.Create(roleUsers).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		services.Logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	for _, user := range role.Users {
		services.RedisStorage.DeleteByMethod(
			fmt.Sprintf("fiber:cache:auth:profile:dashboards:%s", user.UserID),
			"GET",
		)
	}

	for _, user := range twitchUsers.Data.Users {
		services.RedisStorage.DeleteByMethod(
			fmt.Sprintf("fiber:cache:auth:profile:dashboards:%s", user.ID),
			"GET",
		)
	}

	return nil
}
