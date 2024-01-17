package kappagen

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/websockets"
)

var Kappagen = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "kappagen",
		Description: null.StringFrom("Magic ball will answer to all your questions!"),
		Module:      "GAMES",
		IsReply:     true,
		Visible:     true,
		RolesIDS:    pq.StringArray{},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		var emotes []*websockets.TriggerKappagenRequest_Emote
		for _, e := range parseCtx.Emotes {
			emote := &websockets.TriggerKappagenRequest_Emote{
				Id:        e.ID,
				Positions: []string{},
			}

			for _, pos := range e.Positions {
				emote.Positions = append(emote.Positions, fmt.Sprintf("%v-%v", pos.Start, pos.End))
			}

			emotes = append(emotes, emote)
		}

		param := "!" + parseCtx.RawText

		_, err := parseCtx.Services.GrpcClients.WebSockets.TriggerKappagen(
			ctx, &websockets.TriggerKappagenRequest{
				ChannelId: parseCtx.Channel.ID,
				Text:      param,
				Emotes:    emotes,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "error while triggering kappagen",
				Err:     err,
			}
		}

		return &types.CommandsHandlerResult{}, nil
	},
}
