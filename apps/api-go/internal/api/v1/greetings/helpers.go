package greetings

import (
	model "tsuwari/models"
	"tsuwari/twitch"

	"github.com/satont/go-helix/v2"
	"gorm.io/gorm"
)

func getTwitchUserByName(userName string, twitch *twitch.Twitch) *helix.User {
	twitchUsers, err := twitch.Client.GetUsers(&helix.UsersParams{
		Logins: []string{userName},
	})

	if err != nil || len(twitchUsers.Data.Users) == 0 {
		return nil
	}

	twitchUser := twitchUsers.Data.Users[0]
	return &twitchUser
}

func findGreetingByUser(userId string, channelId string, db *gorm.DB) *model.ChannelsGreetings {
	greeting := &model.ChannelsGreetings{}
	err := db.Where(`"channelId" = ? AND "userId" = ?`, channelId, userId).First(greeting).Error
	if err != nil {
		return nil
	}

	return greeting
}
