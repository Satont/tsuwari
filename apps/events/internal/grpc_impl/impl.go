package grpc_impl

import (
	"context"
	"github.com/satont/tsuwari/apps/events/internal"
	"github.com/satont/tsuwari/libs/grpc/generated/events"
	"google.golang.org/protobuf/types/known/emptypb"
)

type EventsGrpcImplementation struct {
	events.UnimplementedEventsServer

	services *internal.Services
}

func NewEvents(services *internal.Services) *EventsGrpcImplementation {
	return &EventsGrpcImplementation{
		services: services,
	}
}

func (c *EventsGrpcImplementation) Follow(_ context.Context, msg *events.FollowMessage) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			UserName:        msg.UserName,
			UserDisplayName: msg.UserDisplayName,
			UserID:          msg.UserId,
		},
		"FOLLOW",
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) Subscribe(_ context.Context, msg *events.SubscribeMessage) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			UserName:        msg.UserName,
			UserDisplayName: msg.UserDisplayName,
			SubLevel:        msg.Level,
			UserID:          msg.UserId,
		},
		"SUBSCRIBE",
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) ReSubscribe(_ context.Context, msg *events.ReSubscribeMessage) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			UserName:        msg.UserName,
			UserDisplayName: msg.UserDisplayName,
			SubLevel:        msg.Level,
			ResubMessage:    msg.Message,
			ResubMonths:     msg.Months,
			ResubStreak:     msg.Streak,
			UserID:          msg.UserId,
		},
		"RESUBSCRIBE",
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) RedemptionCreated(_ context.Context, msg *events.RedemptionCreatedMessage) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			UserName:        msg.UserName,
			UserDisplayName: msg.UserDisplayName,
			RewardCost:      msg.RewardCost,
			RewardInput:     msg.Input,
			RewardName:      msg.RewardName,
			RewardID:        msg.Id,
			UserID:          msg.UserId,
		},
		"REDEMPTION_CREATED",
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) CommandUsed(_ context.Context, msg *events.CommandUsedMessage) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			UserName:        msg.UserName,
			UserDisplayName: msg.UserDisplayName,
			CommandName:     msg.CommandName,
			CommandID:       msg.CommandId,
			CommandInput:    msg.CommandInput,
			UserID:          msg.UserId,
		},
		"COMMAND_USED",
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) FirstUserMessage(_ context.Context, msg *events.FirstUserMessageMessage) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			UserName:        msg.UserName,
			UserDisplayName: msg.UserDisplayName,
			UserID:          msg.UserId,
		},
		"FIRST_USER_MESSAGE",
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) Raided(_ context.Context, msg *events.RaidedMessage) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			UserName:        msg.UserName,
			UserDisplayName: msg.UserDisplayName,
			RaidViewers:     msg.Viewers,
			UserID:          msg.UserId,
		},
		"RAIDED",
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) TitleOrCategoryChanged(_ context.Context, msg *events.TitleOrCategoryChangedMessage) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			OldStreamCategory: msg.OldCategory,
			NewStreamCategory: msg.NewCategory,
			OldStreamTitle:    msg.OldTitle,
			NewStreamTitle:    msg.NewTitle,
		},
		"TITLE_OR_CATEGORY_CHANGED",
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) StreamOnline(_ context.Context, msg *events.StreamOnlineMessage) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			StreamTitle:    msg.Title,
			StreamCategory: msg.Category,
		},
		"STREAM_ONLINE",
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) StreamOffline(_ context.Context, msg *events.StreamOfflineMessage) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{},
		"STREAM_OFFLINE",
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) SubGift(_ context.Context, msg *events.SubGiftMessage) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			TargetUserName:        msg.TargetUserName,
			TargetUserDisplayName: msg.TargetDisplayName,
			SubLevel:              msg.Level,
			UserID:                msg.SenderUserId,
		},
		"SUB_GIFT",
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) ChatClear(_ context.Context, msg *events.ChatClearMessage) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{},
		"ON_CHAT_CLEAR",
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) Donate(_ context.Context, msg *events.DonateMessage) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			UserName:       msg.UserName,
			DonateAmount:   msg.Amount,
			DonateCurrency: msg.Currency,
			DonateMessage:  msg.Message,
		},
		"DONATE",
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) KeywordMatched(
	_ context.Context,
	msg *events.KeywordMatchedMessage,
) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			UserName:        msg.UserName,
			UserDisplayName: msg.UserDisplayName,
			KeywordName:     msg.KeywordName,
			KeywordResponse: msg.KeywordResponse,
			KeywordID:       msg.KeywordId,
			UserID:          msg.UserId,
		},
		"KEYWORD_MATCHED",
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) GreetingSended(_ context.Context, msg *events.GreetingSendedMessage) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			UserName:        msg.UserName,
			UserDisplayName: msg.UserDisplayName,
			UserID:          msg.UserId,
			GreetingText:    msg.GreetingText,
		},
		"GREETING_SENDED",
	)

	return &emptypb.Empty{}, nil
}
