package community

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/twitch"
)

type User struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	DisplayName       string `json:"displayName"`
	Watched           int64  `json:"watched"`
	Messages          int32  `json:"messages"`
	Emotes            int    `json:"emotes"`
	AvatarUrl         string `json:"avatarUrl"`
	UsedChannelPoints string `json:"usedChannelPoints"`
}

func handleGet(channelId, limit, page, sortBy, order string, services *types.Services) ([]User, error) {
	twitchClient, err := twitch.NewAppClient(*services.Config, services.Grpc.Tokens)
	if err != nil {
		services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	parsedLimit, err := strconv.Atoi(limit)
	if err != nil || parsedLimit > 100 {
		return nil, fiber.NewError(http.StatusBadRequest, "invalid limit")
	}

	parsedPage, err := strconv.Atoi(page)
	if err != nil {
		return nil, fiber.NewError(http.StatusBadRequest, "invalid page")
	}

	if sortBy != "watched" && sortBy != "messages" && sortBy != "emotes" && sortBy != "usedChannelPoints" {
		return nil, fiber.NewError(http.StatusBadRequest, "invalid sortBy")
	}

	if order != "desc" && order != "asc" {
		return nil, fiber.NewError(http.StatusBadRequest, "invalid order")
	}

	offset := (parsedPage - 1) * parsedLimit

	channel := &model.Channels{}
	err = services.Gorm.Where(`"id" = ?`, channelId).Find(&channel).Error
	if err != nil {
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}
	if channel.ID == "" {
		return nil, fiber.NewError(http.StatusNotFound, "channel not found")
	}

	query, args, err := squirrel.
		Select(`users_stats.*, COUNT("channels_emotes_usages"."id") AS "emotes"`).
		From("users_stats").
		LeftJoin(`"channels_emotes_usages" ON "channels_emotes_usages"."userId" = "users_stats"."userId" AND "channels_emotes_usages"."channelId" = "users_stats"."channelId"`).
		Where(squirrel.And{
			squirrel.Eq{`"users_stats"."channelId"`: channelId},
			// squirrel.NotEq{`"users_stats"."userId"`: channelId},
			squirrel.NotEq{`"users_stats"."userId"`: channel.BotID},
		}).
		Where(`NOT EXISTS (select 1 from "users_ignored" where "id" = "users_stats"."userId")`).
		Limit(uint64(parsedLimit)).
		Offset(uint64(offset)).
		GroupBy(`"users_stats"."id"`).
		OrderBy(fmt.Sprintf(`"%s" %s`, sortBy, order)).
		ToSql()
	query = services.Sqlx.Rebind(query)

	if err != nil {
		services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	dbUsers := []model.UsersStats{}

	err = services.Sqlx.Select(&dbUsers, query, args...)
	if err != nil {
		services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	twitchUsers, err := twitchClient.GetUsers(&helix.UsersParams{
		IDs: lo.Map(dbUsers, func(record model.UsersStats, _ int) string {
			return record.UserID
		}),
	})
	if err != nil {
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	users := []User{}

	for _, dbUser := range dbUsers {
		twitchUser, ok := lo.Find(twitchUsers.Data.Users, func(item helix.User) bool {
			return item.ID == dbUser.UserID
		})
		if !ok {
			continue
		}

		users = append(users, User{
			ID:                twitchUser.ID,
			Name:              twitchUser.Login,
			DisplayName:       twitchUser.DisplayName,
			Watched:           dbUser.Watched,
			Messages:          dbUser.Messages,
			Emotes:            dbUser.Emotes,
			AvatarUrl:         twitchUser.ProfileImageURL,
			UsedChannelPoints: strconv.FormatInt(dbUser.UsedChannelPoints, 10),
		})
	}

	return users, nil
}

func handleResetStats(channelId string, dto *resetStatsDto, services *types.Services) error {
	var query string
	var args []any
	var sqlQueryErr error

	if dto.Field == "messages" || dto.Field == "watched" || dto.Field == "usedChannelPoints" {
		query, args, sqlQueryErr = squirrel.
			Update("users_stats").
			Where(squirrel.Eq{`"channelId"`: channelId}).
			Set(dto.Field, 0).
			ToSql()
	} else {
		query, args, sqlQueryErr = squirrel.
			Delete("channels_emotes_usages").
			Where(squirrel.Eq{`"channelId"`: channelId}).
			ToSql()
	}

	query = services.Sqlx.Rebind(query)

	if sqlQueryErr != nil {
		services.Logger.Error(sqlQueryErr)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	_, err := services.Sqlx.Query(query, args...)
	if err != nil {
		fmt.Println(query, args)
		services.Logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return nil
}
