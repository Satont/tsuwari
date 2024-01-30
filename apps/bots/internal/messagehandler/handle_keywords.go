package messagehandler

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/lib/pq"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/bots/internal/twitchactions"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/websockets"
)

func (c *MessageHandler) handleKeywords(ctx context.Context, msg handleMessage) error {
	var keywords []model.ChannelsKeywords
	err := c.gorm.WithContext(ctx).Where(
		`"channelId" = ? AND "enabled" = ?`, msg.GetBroadcasterUserId(),
		true,
	).Find(&keywords).Error
	if err != nil {
		return err
	}

	if len(keywords) == 0 {
		return nil
	}

	wg := sync.WaitGroup{}
	wg.Add(len(keywords))

	message := msg.GetMessage().GetText()
	var messagesForSend []string

	matchedKeywords := make([]model.ChannelsKeywords, 0, len(keywords))

	timesInMessage := map[string]int{}

	for _, k := range keywords {
		if k.IsRegular {
			regx, err := regexp.Compile(k.Text)
			if err != nil {
				messagesForSend = append(
					messagesForSend,
					fmt.Sprintf("regular expression is wrong for keyword %s", k.Text),
				)
				continue
			}

			if !regx.MatchString(message) {
				continue
			} else {
				timesInMessage[k.ID] = len(regx.FindAllStringSubmatch(message, -1))
			}
		} else {
			if !strings.Contains(strings.ToLower(message), strings.ToLower(k.Text)) {
				continue
			} else {
				timesInMessage[k.ID] = strings.Count(strings.ToLower(message), strings.ToLower(k.Text))
			}
		}

		isOnCooldown := false
		if k.Cooldown != 0 && k.CooldownExpireAt.Valid {
			isOnCooldown = k.CooldownExpireAt.Time.After(time.Now().UTC())
		}

		if isOnCooldown {
			continue
		}

		matchedKeywords = append(matchedKeywords, k)
	}

	for _, k := range matchedKeywords {
		response := c.keywordsParseResponse(ctx, msg, k)

		c.keywordsTriggerEvent(ctx, msg, k, response)
		c.twitchActions.SendMessage(
			ctx, twitchactions.SendMessageOpts{
				BroadcasterID:        msg.GetBroadcasterUserId(),
				SenderID:             msg.DbChannel.BotID,
				Message:              response,
				ReplyParentMessageID: lo.If(k.IsReply, msg.GetMessageId()).Else(""),
			},
		)
		c.keywordsIncrementStats(ctx, k, timesInMessage[k.ID])
		c.keywordsTriggerAlert(ctx, k)
	}

	wg.Wait()
	return nil
}

func (c *MessageHandler) keywordsIncrementStats(
	ctx context.Context,
	keyword model.ChannelsKeywords,
	count int,
) {
	query := make(map[string]any)
	query["cooldownExpireAt"] = time.Now().
		Add(time.Duration(keyword.Cooldown) * time.Second).
		UTC()

	query["usages"] = keyword.Usages + count
	err := c.gorm.WithContext(ctx).Model(&keyword).Where(
		"id = ?",
		keyword.ID,
	).Select("*").Updates(query).
		Error
	if err != nil {
		c.logger.Error(
			"cannot update keyword usages",
			slog.Any("err", err),
			slog.String("channelId", keyword.ChannelID),
		)
	}
}

func (c *MessageHandler) keywordsTriggerEvent(
	ctx context.Context, msg handleMessage,
	keyword model.ChannelsKeywords, response string,
) {
	_, err := c.eventsGrpc.KeywordMatched(
		ctx,
		&events.KeywordMatchedMessage{
			BaseInfo:        &events.BaseInfo{ChannelId: msg.GetBroadcasterUserId()},
			KeywordId:       keyword.ID,
			KeywordName:     keyword.Text,
			KeywordResponse: response,
			UserId:          msg.GetChatterUserId(),
			UserName:        msg.GetChatterUserLogin(),
			UserDisplayName: msg.GetChatterUserName(),
		},
	)
	if err != nil {
		c.logger.Error(
			"cannot send keywords matched event",
			slog.Any("err", err),
			slog.String("channelId", msg.GetBroadcasterUserId()),
			slog.String("userId", msg.GetChatterUserId()),
		)
	}
}

func (c *MessageHandler) keywordsParseResponse(
	ctx context.Context,
	msg handleMessage,
	keyword model.ChannelsKeywords,
) string {
	if keyword.Response == "" {
		return ""
	}

	requestStruct := &parser.ParseTextRequestData{
		Channel: &parser.Channel{
			Id:   msg.GetMessageId(),
			Name: msg.GetBroadcasterUserLogin(),
		},
		Sender: &parser.Sender{
			Id:          msg.GetChatterUserId(),
			Name:        msg.GetChatterUserLogin(),
			DisplayName: msg.GetChatterUserName(),
			Badges:      createUserBadges(msg.Badges),
		},
		Message: &parser.Message{
			Text: keyword.Response,
			Id:   msg.GetMessageId(),
		},
		ParseVariables: lo.ToPtr(true),
	}

	res, err := c.parserGrpc.ParseTextResponse(ctx, requestStruct)
	if err != nil {
		c.logger.Error(
			"cannot parse keyword response",
			slog.Any("err", err),
			slog.String("channelId", msg.GetBroadcasterUserId()),
		)
	}

	return strings.Join(res.GetResponses(), " ")
}

func (c *MessageHandler) keywordsTriggerAlert(
	ctx context.Context,
	keyword model.ChannelsKeywords,
) {
	alert := model.ChannelAlert{}
	if err := c.gorm.WithContext(ctx).Where(
		"channel_id = ? AND keywords_ids && ?",
		keyword.ChannelID,
		pq.StringArray{keyword.ID},
	).Find(&alert).Error; err != nil {
		c.logger.Error(
			"cannot get alert",
			slog.Any("err", err),
			slog.String("channelId", keyword.ChannelID),
		)
		return
	}

	if alert.ID == "" {
		return
	}

	if _, err := c.websocketsGrpc.TriggerAlert(
		context.Background(),
		&websockets.TriggerAlertRequest{
			ChannelId: keyword.ChannelID,
			AlertId:   alert.ID,
		},
	); err != nil {
		c.logger.Error(
			"cannot trigger alert",
			slog.Any("err", err),
			slog.String("channelId", keyword.ChannelID),
		)
	}
}
