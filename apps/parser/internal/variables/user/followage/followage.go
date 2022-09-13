package userfollowage

import (
	types "tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"
	"tsuwari/parser/pkg/helpers"

	"github.com/nicklaw5/helix"
	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "user.followage",
	Description: lo.ToPtr("User followage"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		targetId := ctx.SenderId
		if ctx.Text != nil {
			users, err := ctx.Services.Twitch.Client.GetUsers(&helix.UsersParams{
				Logins: []string{*ctx.Text},
			})

			if err != nil || len(users.Data.Users) == 0 {
				result.Result = "Cannot find user " + *ctx.Text + " on twitch."
				return result, nil
			}

			targetId = users.Data.Users[0].ID
		}

		if ctx.ChannelId == targetId {
			result.Result = "you cannot fetch followage of yourself because you are broadcaster."
			return result, nil
		}

		follow := ctx.GetFollowAge(targetId)
		if follow == nil {
			result.Result = "not a follower"
		} else {
			result.Result = helpers.Duration(follow.FollowedAt)
		}

		return result, nil
	},
}
