package model

import (
	"github.com/guregu/null"
)

type ChannelModulesSettings struct {
	ID        string      `gorm:"column:id;type:uuid"        json:"id"`
	Type      string      `gorm:"column:type;"               json:"type"`
	Settings  []byte      `gorm:"column:settings;type:jsonb" json:"settings"`
	ChannelId string      `gorm:"column:channelId;type:text" json:"channelId"`
	UserId    null.String `gorm:"column:userId;type:text"    json:"userId"`
}

func (ChannelModulesSettings) TableName() string {
	return "channels_modules_settings"
}

type ChatAlertsSettings struct {
	Followers           ChatAlertsFollowersSettings   `json:"followers"`
	Raids               ChatAlertsRaids               `json:"raids"`
	Donations           ChatAlertsDonations           `json:"donations"`
	Subscribers         ChatAlertsSubscribers         `json:"subscribers"`
	Cheers              ChatAlertsCheers              `json:"cheers"`
	Redemptions         ChatAlertsRedemptions         `json:"redemptions"`
	FirstUserMessage    ChatAlertsFirstUserMessage    `json:"firstUserMessage"`
	StreamOnline        ChatAlertsStreamOnline        `json:"streamOnline"`
	StreamOffline       ChatAlertsStreamOffline       `json:"streamOffline"`
	ChatCleared         ChatAlertsChatCleared         `json:"chatCleared"`
	Ban                 ChatAlertsBan                 `json:"ban"`
	UnbanRequestCreate  ChatAlertsUnbanRequestCreate  `json:"unbanRequestCreate"`
	UnbanRequestResolve ChatAlertsUnbanRequestResolve `json:"unbanRequestResolve"`
	MessageDelete       ChatAlertsMessageDelete       `json:"messageDelete"`
}

type ChatAlertsFollowersSettings struct {
	Enabled  bool                `json:"enabled"`
	Messages []ChatAlertsMessage `json:"messages"`
	Cooldown int                 `json:"cooldown"`
}

type ChatAlertsCountedMessage struct {
	Count int    `json:"count"`
	Text  string `json:"text"`
}

type ChatAlertsMessage struct {
	Text string `json:"text"`
}

type ChatAlertsRaids struct {
	Enabled  bool                       `json:"enabled"`
	Messages []ChatAlertsCountedMessage `json:"messages"`
	Cooldown int                        `json:"cooldown"`
}

type ChatAlertsDonations struct {
	Enabled  bool                       `json:"enabled"`
	Messages []ChatAlertsCountedMessage `json:"messages"`
	Cooldown int                        `json:"cooldown"`
}

type ChatAlertsSubscribers struct {
	Enabled  bool                       `json:"enabled"`
	Messages []ChatAlertsCountedMessage `json:"messages"`
	Cooldown int                        `json:"cooldown"`
}

type ChatAlertsCheers struct {
	Enabled  bool                       `json:"enabled"`
	Messages []ChatAlertsCountedMessage `json:"messages"`
	Cooldown int                        `json:"cooldown"`
}

type ChatAlertsRedemptions struct {
	Enabled  bool                `json:"enabled"`
	Messages []ChatAlertsMessage `json:"messages"`
	Cooldown int                 `json:"cooldown"`
}

type ChatAlertsFirstUserMessage struct {
	Enabled  bool                `json:"enabled"`
	Messages []ChatAlertsMessage `json:"messages"`
	Cooldown int                 `json:"cooldown"`
}

type ChatAlertsStreamOnline struct {
	Enabled  bool                `json:"enabled"`
	Messages []ChatAlertsMessage `json:"messages"`
	Cooldown int                 `json:"cooldown"`
}

type ChatAlertsStreamOffline struct {
	Enabled  bool                `json:"enabled"`
	Messages []ChatAlertsMessage `json:"messages"`
	Cooldown int                 `json:"cooldown"`
}

type ChatAlertsChatCleared struct {
	Enabled  bool                `json:"enabled"`
	Messages []ChatAlertsMessage `json:"messages"`
	Cooldown int                 `json:"cooldown"`
}

type ChatAlertsBan struct {
	Enabled           bool                       `json:"enabled"`
	Messages          []ChatAlertsCountedMessage `json:"messages"`
	IgnoreTimeoutFrom []string                   `json:"ignoreTimeoutFrom"`
	Cooldown          int                        `json:"cooldown"`
}

type ChatAlertsUnbanRequestCreate struct {
	Enabled  bool                `json:"enabled"`
	Messages []ChatAlertsMessage `json:"messages"`
	Cooldown int                 `json:"cooldown"`
}

type ChatAlertsUnbanRequestResolve struct {
	Enabled  bool                `json:"enabled"`
	Messages []ChatAlertsMessage `json:"messages"`
	Cooldown int                 `json:"cooldown"`
}

type ChatAlertsMessageDelete struct {
	Enabled  bool                `json:"enabled"`
	Messages []ChatAlertsMessage `json:"messages"`
	Cooldown int                 `json:"cooldown"`
}
