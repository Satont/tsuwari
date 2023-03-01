package channel_game

import (
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	model "github.com/satont/tsuwari/libs/gomodels"
	"gorm.io/gorm"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

var History = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "game history",
		Description: null.StringFrom("Print history of games."),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MODERATION",
		IsReply:     true,
		Visible:     true,
	},
	Handler: func(ctx *variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		db := do.MustInvoke[gorm.DB](di.Provider)

		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		limit := 5

		if ctx.Text != nil && *ctx.Text != "" {
			l, err := strconv.Atoi(*ctx.Text)
			if err == nil {
				limit = l
			}
		}

		if limit > 20 {
			limit = 5
		}

		histories := []model.ChannelInfoHistory{}
		err := db.
			Raw(`SELECT * FROM (
				SELECT DISTINCT ON (category) * FROM "channels_info_history" 
				                             WHERE "channelId" = ?
				                             ORDER BY "category", "createdAt" 
				                             DESC
				) subquery ORDER BY "createdAt" DESC LIMIT ?`, ctx.ChannelId, limit).
			Find(&histories).
			Error

		if err != nil {
			result.Result = append(result.Result, "internal error")
			return result
		}

		categories := lo.Map(histories, func(item model.ChannelInfoHistory, _ int) string {
			return item.Category
		})

		result.Result = append(result.Result, strings.Join(categories, " ║ "))
		return result
	},
}
