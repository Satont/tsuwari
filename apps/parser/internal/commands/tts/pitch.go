package tts

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	"strconv"
)

var PitchCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "tts pitch",
		Description: lo.ToPtr("Change tts pitch"),
		Permission:  "VIEWER",
		Visible:     true,
		Module:      lo.ToPtr("TTS"),
		IsReply:     true,
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{}
		channelSettings, channelModele := getSettings(ctx.ChannelId, "")

		if channelSettings == nil {
			result.Result = append(result.Result, "TTS is not configured.")
			return result
		}

		userSettings, currentUserModel := getSettings(ctx.ChannelId, ctx.SenderId)

		if ctx.Text == nil {
			result.Result = append(
				result.Result,
				fmt.Sprintf(
					"Global pitch: %v | Your pitch: %v",
					channelSettings.Pitch,
					lo.IfF(userSettings != nil, func() int {
						return userSettings.Pitch
					}).Else(channelSettings.Pitch),
				))
			return result
		}

		pitch, err := strconv.Atoi(*ctx.Text)
		if err != nil {
			result.Result = append(result.Result, "Pitch must be a number")
			return result
		}

		if pitch < 0 || pitch > 100 {
			result.Result = append(result.Result, "Pitch must be between 0 and 100")
			return result
		}

		if ctx.ChannelId == ctx.SenderId {
			channelSettings.Pitch = pitch
			err := updateSettings(channelModele, channelSettings)
			if err != nil {
				fmt.Println(err)
				result.Result = append(result.Result, "Error while updating settings")
				return result
			}
		} else {
			if userSettings == nil {
				_, _, err := createUserSettings(pitch, 50, channelSettings.Voice, ctx.ChannelId, ctx.SenderId)
				if err != nil {
					fmt.Println(err)
					result.Result = append(result.Result, "Error while creating settings")
					return result
				}
			} else {
				userSettings.Pitch = pitch
				err := updateSettings(currentUserModel, userSettings)
				if err != nil {
					fmt.Println(err)
					result.Result = append(result.Result, "Error while updating settings")
					return result
				}
			}
		}

		result.Result = append(result.Result, fmt.Sprintf("Pitch changed to %v", pitch))

		return result
	},
}
