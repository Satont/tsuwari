package chat_client

import (
	"context"
	"log/slog"

	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/grpc/parser"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *ChatClient) handleCommand(msg Message, userBadges []string) {
	requestStruct := &parser.ProcessCommandRequest{
		Sender: &parser.Sender{
			Id:          msg.User.ID,
			Name:        msg.User.Name,
			DisplayName: msg.User.DisplayName,
			Badges:      userBadges,
		},
		Channel: &parser.Channel{
			Id:   msg.Channel.ID,
			Name: msg.Channel.Name,
		},
		Message: &parser.Message{
			Id:   msg.ID,
			Text: msg.Message,
			Emotes: lo.Map(
				msg.Emotes,
				func(item *irc.Emote, _ int) *parser.Message_Emote {
					return &parser.Message_Emote{
						Name:  item.Name,
						Id:    item.ID,
						Count: int64(item.Count),
						Positions: lo.Map(
							item.Positions,
							func(item irc.EmotePosition, _ int) *parser.Message_EmotePosition {
								return &parser.Message_EmotePosition{
									Start: int64(item.Start),
									End:   int64(item.End),
								}
							},
						),
					}
				},
			),
		},
	}

	res, err := c.services.ParserGrpc.ProcessCommand(context.Background(), requestStruct)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return
		}

		c.services.Logger.Error("cannot process command", slog.Any("err", err))
		return
	}

	if res.KeepOrder != nil && !*res.KeepOrder {
		for _, r := range res.Responses {
			if r == "" || r == " " {
				continue
			}

			c.Say(
				SayOpts{
					Channel:   msg.Channel.Name,
					Text:      r,
					ReplyTo:   lo.If(res.IsReply, lo.ToPtr(msg.ID)).Else(nil),
					WithLimit: true,
				},
			)
		}
	} else {
		for _, r := range res.Responses {
			if r == "" || r == " " {
				continue
			}

			go c.Say(
				SayOpts{
					Channel:   msg.Channel.Name,
					Text:      r,
					ReplyTo:   lo.If(res.IsReply, lo.ToPtr(msg.ID)).Else(nil),
					WithLimit: true,
				},
			)
		}
	}
}
