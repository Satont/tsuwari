package messages

import (
	model "github.com/satont/twir/libs/gomodels"
	"go.uber.org/zap"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func IncrementUserMessages(db *gorm.DB, userId, channelId string) {
	stream := model.ChannelsStreams{}
	if err := db.Where(`"userId" = ?`, channelId).Find(&stream).Error; err != nil {
		zap.S().Error(err)
		return
	}

	if stream.ID == "" {
		return
	}

	user := model.Users{}
	err := db.
		Where(`"id" = ?`, userId).
		Preload("Stats", `"userId" = ? AND "channelId" = ?`, userId, channelId).
		Find(&user).Error
	if err != nil {
		zap.S().Error(err)
		return
	}

	// no user found
	if user.ID == "" {
		user.ID = userId
		user.ApiKey = uuid.NewV4().String()
		user.IsBotAdmin = false
		user.IsTester = false
		user.Stats = createStats(userId, channelId)

		if err := db.Create(&user).Error; err != nil {
			zap.S().Error(err)
			return
		}
	} else {
		if user.Stats == nil {
			newStats := createStats(userId, channelId)
			err := db.Create(newStats).Error
			if err != nil {
				zap.S().Error(err, newStats)
			}
		} else {
			err := db.
				Model(&user.Stats).
				Where(`"userId" = ? AND "channelId" = ?`, userId, channelId).
				Update("messages", user.Stats.Messages+1).
				Error
			if err != nil {
				zap.S().Error(err)
			}
		}
	}
}

func createStats(userId, channelId string) *model.UsersStats {
	stats := &model.UsersStats{
		ID:                uuid.NewV4().String(),
		UserID:            userId,
		ChannelID:         channelId,
		Messages:          1,
		Watched:           0,
		UsedChannelPoints: 0,
	}
	return stats
}
