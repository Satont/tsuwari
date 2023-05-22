package tts

import (
	"context"

	"github.com/guregu/null"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
)

var DisableCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "tts disable",
		Description: null.StringFrom("Disable tts."),
		Module:      "TTS",
		IsReply:     true,
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{}
		err := switchEnableState(ctx, parseCtx.Services.Gorm, parseCtx.Channel.ID, false)
		if err != nil {
			result.Result = append(result.Result, "Error while updating settings")
			return result
		}

		result.Result = append(result.Result, "TTS disabled")

		return result
	},
}
