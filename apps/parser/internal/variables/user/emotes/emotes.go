package user_emotes

import (
	"fmt"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	model "github.com/satont/tsuwari/libs/gomodels"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var Variable = types.Variable{
	Name:         "user.emotes",
	Description:  lo.ToPtr("User used emotes count"),
	CommandsOnly: lo.ToPtr(true),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		db := do.MustInvoke[gorm.DB](di.Provider)
		result := &types.VariableHandlerResult{}

		var count int64
		err := db.
			Where(`"channelId" = ? AND "userId" = ?`, ctx.ChannelId, ctx.SenderId).
			Model(&model.ChannelEmoteUsage{}).
			Count(&count).
			Error

		if err != nil {
			zap.S().Error(err)
			result.Result = "error"
			return result, nil
		}

		result.Result = fmt.Sprint(count)

		return result, nil
	},
}
