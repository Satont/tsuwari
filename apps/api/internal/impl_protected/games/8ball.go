package games

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/games"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

var eightBallType = "8ball"

func (c *Games) GamesGetEightBallSettings(
	ctx context.Context,
	_ *emptypb.Empty,
) (*games.EightBallSettingsResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	entity := model.ChannelModulesSettings{}
	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? and "userId" IS NULL and type = ?`, dashboardId, eightBallType).
		First(&entity).
		Error; err != nil {
		return nil, err
	}

	settings := model.EightBallSettings{}
	if err := json.Unmarshal(entity.Settings, &settings); err != nil {
		return nil, err
	}

	return &games.EightBallSettingsResponse{
		Answers: settings.Answers,
		Enabled: settings.Enabled,
	}, nil
}

const maxAnswers = 25

func (c *Games) GamesUpdateEightBallSettings(
	ctx context.Context,
	req *games.UpdateEightBallSettings,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	if len(req.Answers) > maxAnswers {
		return nil, twirp.NewError("400", fmt.Sprintf("Max answers is %v", maxAnswers))
	}

	entity := model.ChannelModulesSettings{}
	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? and "userId" IS NULL and type = ?`, dashboardId, eightBallType).
		Find(&entity).
		Error; err != nil {
		return nil, err
	}

	if entity.ID == "" {
		entity.ID = uuid.New().String()
		entity.ChannelId = dashboardId
		entity.Type = eightBallType
	}

	settings := model.EightBallSettings{
		Answers: req.Answers,
		Enabled: req.Enabled,
	}

	settingsJson, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	entity.Settings = settingsJson

	txErr := c.Db.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			if err := tx.WithContext(ctx).Save(&entity).Error; err != nil {
				return err
			}

			eightBallCommand := model.ChannelsCommands{}
			if err := tx.WithContext(ctx).Where(
				`"channelId" = ? and "defaultName" = ?`,
				dashboardId,
				"8ball",
			).First(&eightBallCommand).Error; err != nil {
				return err
			}

			eightBallCommand.Enabled = req.Enabled

			if err := tx.WithContext(ctx).Save(&eightBallCommand).Error; err != nil {
				return err
			}

			return nil
		},
	)

	if txErr != nil {
		return nil, fmt.Errorf("transaction error: %w", txErr)
	}

	return &emptypb.Empty{}, nil
}
