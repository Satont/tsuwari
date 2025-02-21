package resolvers

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger/audit"
	"github.com/satont/twir/libs/types/types/api/overlays"
	"github.com/satont/twir/libs/utils"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	now_playing_fetcher "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/now-playing-fetcher"
)

func chatOverlayDbToGql(entity *model.ChatOverlaySettings) *gqlmodel.ChatOverlay {
	return &gqlmodel.ChatOverlay{
		ID:                  entity.ID.String(),
		MessageHideTimeout:  int(entity.MessageHideTimeout),
		MessageShowDelay:    int(entity.MessageShowDelay),
		Preset:              entity.Preset,
		FontSize:            int(entity.FontSize),
		HideCommands:        entity.HideCommands,
		HideBots:            entity.HideBots,
		FontFamily:          entity.FontFamily,
		ShowBadges:          entity.ShowBadges,
		ShowAnnounceBadge:   entity.ShowAnnounceBadge,
		TextShadowColor:     entity.TextShadowColor,
		TextShadowSize:      int(entity.TextShadowSize),
		ChatBackgroundColor: entity.ChatBackgroundColor,
		Direction:           entity.Direction,
		FontWeight:          int(entity.FontWeight),
		FontStyle:           entity.FontStyle,
		PaddingContainer:    int(entity.PaddingContainer),
		Animation:           gqlmodel.ChatOverlayAnimation(entity.Animation),
	}
}

func (r *queryResolver) chatOverlays(
	ctx context.Context,
) ([]gqlmodel.ChatOverlay, error) {
	dashboardId, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	var entities []model.ChatOverlaySettings
	if err := r.deps.Gorm.
		WithContext(ctx).
		Where("channel_id = ?", dashboardId).
		Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("failed to get chat overlay settings: %w", err)
	}

	var result []gqlmodel.ChatOverlay
	for _, entity := range entities {
		result = append(result, *chatOverlayDbToGql(&entity))
	}

	return result, nil
}

func (r *Resolver) getChatOverlaySettings(
	ctx context.Context,
	id,
	channelId string,
) (*gqlmodel.ChatOverlay, error) {
	entity := model.ChatOverlaySettings{
		ID:        uuid.MustParse(id),
		ChannelID: channelId,
	}
	if err := r.deps.Gorm.
		WithContext(ctx).
		First(&entity).Error; err != nil {
		return nil, fmt.Errorf("failed to get chat overlay settings: %w", err)
	}

	return chatOverlayDbToGql(&entity), nil
}

func (r *mutationResolver) updateChatOverlay(
	ctx context.Context,
	id string,
	opts gqlmodel.ChatOverlayMutateOpts,
) (bool, error) {
	dashboardId, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	user, err := r.deps.Sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	entity := model.ChatOverlaySettings{}
	if err := r.deps.Gorm.
		WithContext(ctx).
		Where("channel_id = ? AND id = ?", dashboardId, id).
		First(&entity).Error; err != nil {
		return false, fmt.Errorf("failed to get chat overlay settings: %w", err)
	}

	var entityCopy model.ChatOverlaySettings
	if err := utils.DeepCopy(&entity, &entityCopy); err != nil {
		return false, fmt.Errorf("failed to copy chat overlay settings: %w", err)
	}

	if opts.MessageHideTimeout.IsSet() {
		entity.MessageHideTimeout = uint32(*opts.MessageHideTimeout.Value())
	}

	if opts.MessageShowDelay.IsSet() {
		entity.MessageShowDelay = uint32(*opts.MessageShowDelay.Value())
	}

	if opts.Preset.IsSet() {
		entity.Preset = *opts.Preset.Value()
	}

	if opts.FontSize.IsSet() {
		entity.FontSize = uint32(*opts.FontSize.Value())
	}

	if opts.HideCommands.IsSet() {
		entity.HideCommands = *opts.HideCommands.Value()
	}

	if opts.HideBots.IsSet() {
		entity.HideBots = *opts.HideBots.Value()
	}

	if opts.FontFamily.IsSet() {
		entity.FontFamily = *opts.FontFamily.Value()
	}

	if opts.ShowBadges.IsSet() {
		entity.ShowBadges = *opts.ShowBadges.Value()
	}

	if opts.ShowAnnounceBadge.IsSet() {
		entity.ShowAnnounceBadge = *opts.ShowAnnounceBadge.Value()
	}

	if opts.TextShadowColor.IsSet() {
		entity.TextShadowColor = *opts.TextShadowColor.Value()
	}

	if opts.TextShadowSize.IsSet() {
		entity.TextShadowSize = uint32(*opts.TextShadowSize.Value())
	}

	if opts.ChatBackgroundColor.IsSet() {
		entity.ChatBackgroundColor = *opts.ChatBackgroundColor.Value()
	}

	if opts.Direction.IsSet() {
		entity.Direction = *opts.Direction.Value()
	}

	if opts.FontWeight.IsSet() {
		entity.FontWeight = uint32(*opts.FontWeight.Value())
	}

	if opts.FontStyle.IsSet() {
		entity.FontStyle = *opts.FontStyle.Value()
	}

	if opts.PaddingContainer.IsSet() {
		entity.PaddingContainer = uint32(*opts.PaddingContainer.Value())
	}

	if opts.Animation.IsSet() {
		entity.Animation = model.ChatOverlaySettingsAnimationType(*opts.Animation.Value())
	}

	if err := r.deps.Gorm.
		WithContext(ctx).
		Save(&entity).Error; err != nil {
		return false, fmt.Errorf("failed to save chat overlay settings: %w", err)
	}

	go func() {
		if err := r.deps.WsRouter.Publish(
			chatOverlaySubscriptionKeyCreate(entity.ID.String(), entity.ChannelID),
			chatOverlayDbToGql(&entity),
		); err != nil {
			r.deps.Logger.Error("failed to publish settings", slog.Any("err", err))
		}
	}()

	r.deps.Logger.Audit(
		"Chat overlay update",
		audit.Fields{
			OldValue:      entityCopy,
			NewValue:      entity,
			ActorID:       lo.ToPtr(user.ID),
			ChannelID:     lo.ToPtr(dashboardId),
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelOverlayChat),
			OperationType: audit.OperationUpdate,
		},
	)

	return true, nil
}

func (r *mutationResolver) chatOverlayCreate(
	ctx context.Context,
	opts gqlmodel.ChatOverlayMutateOpts,
) (bool, error) {
	dashboardId, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	user, err := r.deps.Sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	entity := model.ChatOverlaySettings{
		ID:        uuid.UUID{},
		ChannelID: dashboardId,
	}

	if opts.MessageHideTimeout.IsSet() {
		entity.MessageHideTimeout = uint32(*opts.MessageHideTimeout.Value())
	}

	if opts.MessageShowDelay.IsSet() {
		entity.MessageShowDelay = uint32(*opts.MessageShowDelay.Value())
	}

	if opts.Preset.IsSet() {
		entity.Preset = *opts.Preset.Value()
	}

	if opts.FontSize.IsSet() {
		entity.FontSize = uint32(*opts.FontSize.Value())
	}

	if opts.HideCommands.IsSet() {
		entity.HideCommands = *opts.HideCommands.Value()
	}

	if opts.HideBots.IsSet() {
		entity.HideBots = *opts.HideBots.Value()
	}

	if opts.FontFamily.IsSet() {
		entity.FontFamily = *opts.FontFamily.Value()
	}

	if opts.ShowBadges.IsSet() {
		entity.ShowBadges = *opts.ShowBadges.Value()
	}

	if opts.ShowAnnounceBadge.IsSet() {
		entity.ShowAnnounceBadge = *opts.ShowAnnounceBadge.Value()
	}

	if opts.TextShadowColor.IsSet() {
		entity.TextShadowColor = *opts.TextShadowColor.Value()
	}

	if opts.TextShadowSize.IsSet() {
		entity.TextShadowSize = uint32(*opts.TextShadowSize.Value())
	}

	if opts.ChatBackgroundColor.IsSet() {
		entity.ChatBackgroundColor = *opts.ChatBackgroundColor.Value()
	}

	if opts.Direction.IsSet() {
		entity.Direction = *opts.Direction.Value()
	}

	if opts.FontWeight.IsSet() {
		entity.FontWeight = uint32(*opts.FontWeight.Value())
	}

	if opts.FontStyle.IsSet() {
		entity.FontStyle = *opts.FontStyle.Value()
	}

	if opts.PaddingContainer.IsSet() {
		entity.PaddingContainer = uint32(*opts.PaddingContainer.Value())
	}

	if opts.Animation.IsSet() {
		entity.Animation = model.ChatOverlaySettingsAnimationType(*opts.Animation.Value())
	}

	if err := r.deps.Gorm.
		WithContext(ctx).
		Create(&entity).Error; err != nil {
		return false, fmt.Errorf("failed to create chat overlay settings: %w", err)
	}

	r.deps.Logger.Audit(
		"Chat overlay create",
		audit.Fields{
			NewValue:      entity,
			ActorID:       lo.ToPtr(user.ID),
			ChannelID:     lo.ToPtr(dashboardId),
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelOverlayChat),
			OperationType: audit.OperationCreate,
		},
	)

	return true, nil
}

func (r *mutationResolver) chatOverlayDelete(ctx context.Context, id string) (bool, error) {
	dashboardId, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	user, err := r.deps.Sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	entity := model.ChatOverlaySettings{}
	if err := r.deps.Gorm.
		WithContext(ctx).
		Where("channel_id = ? AND id = ?", dashboardId, id).
		Delete(&entity).Error; err != nil {
		return false, fmt.Errorf("failed to delete chat overlay settings: %w", err)
	}

	r.deps.Logger.Audit(
		"Chat overlay delete",
		audit.Fields{
			OldValue:      entity,
			NewValue:      nil,
			ActorID:       lo.ToPtr(user.ID),
			ChannelID:     lo.ToPtr(dashboardId),
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelOverlayChat),
			OperationType: audit.OperationDelete,
		},
	)

	return true, nil
}

func (r *queryResolver) nowPlayingOverlays(ctx context.Context) (
	[]gqlmodel.NowPlayingOverlay,
	error,
) {
	dashboardId, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	var entities []model.ChannelOverlayNowPlaying
	if err := r.deps.Gorm.
		WithContext(ctx).
		Where("channel_id = ?", dashboardId).
		Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("failed to get now playing overlay settings: %w", err)
	}

	result := make([]gqlmodel.NowPlayingOverlay, 0, len(entities))
	for _, entity := range entities {
		var hideTimeout *int
		if entity.HideTimeout.Valid {
			hideTimeout = lo.ToPtr(int(*entity.HideTimeout.Ptr()))
		}

		result = append(
			result,
			gqlmodel.NowPlayingOverlay{
				ID:              entity.ID.String(),
				Preset:          gqlmodel.NowPlayingOverlayPreset(entity.Preset.String()),
				ChannelID:       entity.ChannelID,
				FontFamily:      entity.FontFamily,
				FontWeight:      int(entity.FontWeight),
				BackgroundColor: entity.BackgroundColor,
				ShowImage:       entity.ShowImage,
				HideTimeout:     hideTimeout,
			},
		)
	}

	return result, nil
}

func (r *Resolver) getNowPlayingOverlaySettings(ctx context.Context, id, dashboardID string) (
	*gqlmodel.NowPlayingOverlay,
	error,
) {
	entity := model.ChannelOverlayNowPlaying{
		ID:        uuid.MustParse(id),
		ChannelID: dashboardID,
	}
	if err := r.deps.Gorm.
		WithContext(ctx).
		First(&entity).Error; err != nil {
		return nil, fmt.Errorf("failed to get now playing overlay settings: %w", err)
	}

	var hideTimeout *int
	if entity.HideTimeout.Valid {
		hideTimeout = lo.ToPtr(int(*entity.HideTimeout.Ptr()))
	}

	return &gqlmodel.NowPlayingOverlay{
		ID:              entity.ID.String(),
		Preset:          gqlmodel.NowPlayingOverlayPreset(entity.Preset.String()),
		ChannelID:       entity.ChannelID,
		FontFamily:      entity.FontFamily,
		FontWeight:      int(entity.FontWeight),
		BackgroundColor: entity.BackgroundColor,
		ShowImage:       entity.ShowImage,
		HideTimeout:     hideTimeout,
	}, nil
}

func (r *mutationResolver) deleteNowPlayingOverlay(ctx context.Context, id string) (bool, error) {
	dashboardID, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	user, err := r.deps.Sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	entity := model.ChannelOverlayNowPlaying{
		ID:        uuid.MustParse(id),
		ChannelID: dashboardID,
	}
	if err := r.deps.Gorm.
		WithContext(ctx).
		Delete(&entity).Error; err != nil {
		return false, fmt.Errorf("failed to delete now playing overlay settings: %w", err)
	}

	r.deps.Logger.Audit(
		"Now playing overlay delete",
		audit.Fields{
			OldValue:      entity,
			NewValue:      nil,
			ActorID:       lo.ToPtr(user.ID),
			ChannelID:     lo.ToPtr(dashboardID),
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelOverlayNowPlaying),
			OperationType: audit.OperationDelete,
			ObjectID:      lo.ToPtr(entity.ID.String()),
		},
	)

	return true, nil
}

func (r *mutationResolver) createNowPlayingOverlay(
	ctx context.Context,
	opts gqlmodel.NowPlayingOverlayMutateOpts,
) (bool, error) {
	dashboardID, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	user, err := r.deps.Sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	entity := model.ChannelOverlayNowPlaying{
		ID:        uuid.New(),
		ChannelID: dashboardID,
	}

	if opts.Preset.IsSet() {
		entity.Preset = overlays.ChannelOverlayNowPlayingPreset(opts.Preset.Value().String())
	}

	if opts.FontFamily.IsSet() {
		entity.FontFamily = *opts.FontFamily.Value()
	}

	if opts.FontWeight.IsSet() {
		entity.FontWeight = uint32(*opts.FontWeight.Value())
	}

	if opts.BackgroundColor.IsSet() {
		entity.BackgroundColor = *opts.BackgroundColor.Value()
	}

	if opts.ShowImage.IsSet() {
		entity.ShowImage = *opts.ShowImage.Value()
	}

	if opts.HideTimeout.IsSet() {
		if opts.HideTimeout.Value() == nil {
			entity.HideTimeout = null.Int{}
		} else {
			entity.HideTimeout = null.IntFrom(int64(*opts.HideTimeout.Value()))
		}
	}

	if err := r.deps.Gorm.
		WithContext(ctx).
		Create(&entity).Error; err != nil {
		return false, fmt.Errorf("failed to create now playing overlay settings: %w", err)
	}

	r.deps.Logger.Audit(
		"Now playing overlay create",
		audit.Fields{
			NewValue:      entity,
			ActorID:       lo.ToPtr(user.ID),
			ChannelID:     lo.ToPtr(dashboardID),
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelOverlayNowPlaying),
			OperationType: audit.OperationCreate,
		},
	)

	return true, nil
}

func (r *mutationResolver) updateNowPlayingOverlay(
	ctx context.Context,
	id string,
	opts gqlmodel.NowPlayingOverlayMutateOpts,
) (bool, error) {
	dashboardID, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	user, err := r.deps.Sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	entity := model.ChannelOverlayNowPlaying{
		ID:        uuid.MustParse(id),
		ChannelID: dashboardID,
	}
	if err := r.deps.Gorm.
		WithContext(ctx).
		First(&entity).Error; err != nil {
		return false, fmt.Errorf("failed to get now playing overlay settings: %w", err)
	}

	var entityCopy model.ChannelOverlayNowPlaying
	if err := utils.DeepCopy(&entity, &entityCopy); err != nil {
		return false, fmt.Errorf("failed to copy now playing overlay settings: %w", err)
	}

	if opts.Preset.IsSet() {
		entity.Preset = overlays.ChannelOverlayNowPlayingPreset(opts.Preset.Value().String())
	}

	if opts.FontFamily.IsSet() {
		entity.FontFamily = *opts.FontFamily.Value()
	}

	if opts.FontWeight.IsSet() {
		entity.FontWeight = uint32(*opts.FontWeight.Value())
	}

	if opts.BackgroundColor.IsSet() {
		entity.BackgroundColor = *opts.BackgroundColor.Value()
	}

	if opts.ShowImage.IsSet() {
		entity.ShowImage = *opts.ShowImage.Value()
	}

	if opts.HideTimeout.IsSet() {
		if opts.HideTimeout.Value() == nil {
			entity.HideTimeout = null.Int{}
		} else {
			entity.HideTimeout = null.IntFrom(int64(*opts.HideTimeout.Value()))
		}
	}

	if err := r.deps.Gorm.
		WithContext(ctx).
		Save(&entity).Error; err != nil {
		return false, fmt.Errorf("failed to save now playing overlay settings: %w", err)
	}

	go func() {
		if err := r.deps.WsRouter.Publish(
			nowPlayingOverlaySubscriptionKeyCreate(entity.ID.String(), entity.ChannelID),
			&gqlmodel.NowPlayingOverlay{
				ID:              entity.ID.String(),
				Preset:          gqlmodel.NowPlayingOverlayPreset(entity.Preset.String()),
				ChannelID:       entity.ChannelID,
				FontFamily:      entity.FontFamily,
				FontWeight:      int(entity.FontWeight),
				BackgroundColor: entity.BackgroundColor,
				ShowImage:       entity.ShowImage,
				HideTimeout:     lo.ToPtr(int(entity.HideTimeout.Int64)),
			},
		); err != nil {
			r.deps.Logger.Error("failed to publish settings", slog.Any("err", err))
		}
	}()

	r.deps.Logger.Audit(
		"Now playing overlay update",
		audit.Fields{
			OldValue:      entityCopy,
			NewValue:      entity,
			ActorID:       lo.ToPtr(user.ID),
			ChannelID:     lo.ToPtr(dashboardID),
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelOverlayNowPlaying),
			OperationType: audit.OperationUpdate,
			ObjectID:      lo.ToPtr(entity.ID.String()),
		},
	)

	return true, nil
}

func (r *subscriptionResolver) nowPlayingOverlaySettingsSubscription(
	ctx context.Context,
	id string,
	apiKey string,
) (<-chan *gqlmodel.NowPlayingOverlay, error) {
	user := model.Users{}
	if err := r.deps.Gorm.
		WithContext(ctx).
		Where(`"apiKey" = ?`, apiKey).
		First(&user).
		Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	channel := make(chan *gqlmodel.NowPlayingOverlay)

	go func() {
		sub, err := r.deps.WsRouter.Subscribe(
			[]string{
				nowPlayingOverlaySubscriptionKeyCreate(id, user.ID),
			},
		)
		if err != nil {
			panic(err)
		}
		defer func() {
			sub.Unsubscribe()
			close(channel)
		}()

		initialSettings, err := r.getNowPlayingOverlaySettings(ctx, id, user.ID)
		if err == nil {
			channel <- initialSettings
		}

		for {
			select {
			case <-ctx.Done():
				return
			case data := <-sub.GetChannel():
				var settings gqlmodel.NowPlayingOverlay
				if err := json.Unmarshal(data, &settings); err != nil {
					panic(err)
				}

				channel <- &settings
			}
		}
	}()

	return channel, nil
}

func (r *subscriptionResolver) nowPlayingCurrentTrackSubscription(
	ctx context.Context,
	apiKey string,
) (<-chan *gqlmodel.NowPlayingOverlayTrack, error) {
	user := model.Users{}
	if err := r.deps.Gorm.
		WithContext(ctx).
		Where(`"apiKey" = ?`, apiKey).
		First(&user).
		Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	npService, err := now_playing_fetcher.New(
		now_playing_fetcher.Opts{
			Gorm:      r.deps.Gorm,
			ChannelID: user.ID,
			Redis:     r.deps.Redis,
			Logger:    r.deps.Logger,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create now playing service: %w", err)
	}

	channel := make(chan *gqlmodel.NowPlayingOverlayTrack)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				track, err := npService.Fetch(ctx)
				if err != nil {
					r.deps.Logger.Error("failed to get now playing track", slog.Any("err", err))
					time.Sleep(5 * time.Second)
					continue
				}

				if track == nil {
					channel <- nil
					time.Sleep(5 * time.Second)
					continue
				}

				var imageUrl *string
				if track.ImageUrl != "" {
					imageUrl = &track.ImageUrl
				}

				channel <- &gqlmodel.NowPlayingOverlayTrack{
					Artist:     track.Artist,
					Title:      track.Title,
					ImageURL:   imageUrl,
					ProgressMs: track.ProgressMs,
					DurationMs: track.DurationMs,
				}

				time.Sleep(1 * time.Second)
			}
		}
	}()

	return channel, nil
}
