package manage

import (
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"gorm.io/gorm"
	"log"
	"strings"

	"github.com/satont/tsuwari/apps/parser/internal/types"

	model "github.com/satont/tsuwari/libs/gomodels"

	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	"github.com/samber/lo"
)

var RemoveAliaseCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "commands aliases remove",
		Description: null.StringFrom("Remove aliase from command"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MANAGE",
		IsReply:     true,
	},
	Handler: func(ctx *variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		db := do.MustInvoke[gorm.DB](di.Provider)

		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		if ctx.Text == nil {
			result.Result = append(result.Result, incorrectUsage)
			return result
		}

		args := strings.Split(*ctx.Text, " ")

		if len(args) < 2 {
			result.Result = append(result.Result, incorrectUsage)
			return result
		}

		commandName := strings.ToLower(strings.ReplaceAll(args[0], "!", ""))
		aliase := strings.ToLower(strings.ReplaceAll(strings.Join(args[1:], " "), "!", ""))

		cmd := model.ChannelsCommands{}
		err := db.
			Where(`"channelId" = ? AND name = ?`, ctx.ChannelId, commandName).
			First(&cmd).Error

		if err != nil || cmd.ID == "" {
			result.Result = append(result.Result, "Command not found.")
			return result
		}

		if !lo.Contains(cmd.Aliases, aliase) {
			result.Result = append(result.Result, "That aliase not in the command")
			return result
		}

		index := lo.IndexOf(cmd.Aliases, aliase)
		cmd.Aliases = append(cmd.Aliases[:index], cmd.Aliases[index+1:]...)

		err = db.
			Save(&cmd).Error

		if err != nil {
			log.Fatalln(err)
			result.Result = append(
				result.Result,
				"Cannot update command aliases. This is internal bug, please report it.",
			)
			return result
		}

		result.Result = append(result.Result, "✅ Aliase removed.")
		return result
	},
}
