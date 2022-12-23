package greetings

import (
	"net/http"
	"sync"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
	uuid "github.com/satori/go.uuid"
)

type Greeting struct {
	model.ChannelsGreetings
	UserName string `json:"userName"`
	Avatar   string `json:"avatar"`
}

func handleGet(channelId string, services types.Services) []Greeting {
	greetings := []model.ChannelsGreetings{}
	services.DB.Where(`"channelId" = ?`, channelId).Find(&greetings)
	users := []Greeting{}

	chunks := lo.Chunk(greetings, 100)
	greetingsWg := sync.WaitGroup{}
	greetingsWg.Add(len(chunks))
	for _, chunk := range chunks {
		go func(c []model.ChannelsGreetings) {
			defer greetingsWg.Done()
			ids := lo.Map(chunk, func(g model.ChannelsGreetings, _ int) string {
				return g.UserID
			})
			twitchUsers, err := services.Twitch.Client.GetUsers(&helix.UsersParams{
				IDs: ids,
			})
			if err != nil {
				return
			}
			for _, u := range twitchUsers.Data.Users {
				user, ok := lo.Find(greetings, func(g model.ChannelsGreetings) bool {
					return g.UserID == u.ID
				})
				if ok {
					users = append(users, Greeting{
						ChannelsGreetings: user,
						UserName:          u.Login,
						Avatar:            u.ProfileImageURL,
					})
				}
			}
		}(chunk)
	}

	greetingsWg.Wait()

	return users
}

func handlePost(
	channelId string,
	dto *greetingsDto,
	services types.Services,
) (*Greeting, error) {
	twitchUser := getTwitchUserByName(dto.Username, services.Twitch)
	if twitchUser == nil {
		return nil, fiber.NewError(http.StatusNotFound, "cannot find twitch user")
	}

	existedGreeting := findGreetingByUser(twitchUser.ID, channelId, services.DB)
	if existedGreeting != nil {
		return nil, fiber.NewError(400, "greeting for that user already exists")
	}

	greeting := &model.ChannelsGreetings{
		ID:        uuid.NewV4().String(),
		ChannelID: channelId,
		UserID:    twitchUser.ID,
		Enabled:   *dto.Enabled,
		Text:      dto.Text,
		IsReply:   *dto.IsReply,
		Processed: false,
	}

	err := services.DB.Save(greeting).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot create greeting")
	}

	return &Greeting{
		ChannelsGreetings: *greeting,
		UserName:          twitchUser.Login,
	}, nil
}

func handleDelete(greetingId string, services types.Services) error {
	greeting := findGreetingById(greetingId, services.DB)
	if greeting == nil {
		return fiber.NewError(http.StatusNotFound, "greeting not found")
	}
	err := services.DB.Where("id = ?", greetingId).Delete(&model.ChannelsGreetings{}).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot delete greeting")
	}

	return nil
}

func handleUpdate(
	greetingId string,
	dto *greetingsDto,
	services types.Services,
) (*Greeting, error) {
	greeting := findGreetingById(greetingId, services.DB)
	if greeting == nil {
		return nil, fiber.NewError(http.StatusNotFound, "greeting not found")
	}

	newGreeting := &model.ChannelsGreetings{
		ID:        greeting.ID,
		ChannelID: greeting.ChannelID,
		UserID:    greeting.UserID,
		Enabled:   *dto.Enabled,
		Text:      dto.Text,
		IsReply:   *dto.IsReply,
		Processed: false,
	}

	twitchUser := getTwitchUserByName(dto.Username, services.Twitch)
	if twitchUser == nil {
		return nil, fiber.NewError(http.StatusNotFound, "cannot find twitch user")
	}

	newGreeting.UserID = twitchUser.ID

	err := services.DB.Model(greeting).Select("*").Updates(newGreeting).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot update greeting")
	}

	return &Greeting{
		ChannelsGreetings: *newGreeting,
		UserName:          twitchUser.Login,
		Avatar:            twitchUser.ProfileImageURL,
	}, nil
}

func handlePatch(
	channelId, greetingId string,
	dto *greetingsPatchDto,
	services types.Services,
) (*Greeting, error) {
	greeting := findGreetingById(greetingId, services.DB)
	if greeting == nil {
		return nil, fiber.NewError(http.StatusNotFound, "greeting not found")
	}

	twitchUser := getTwitchUserById(greeting.UserID, services.Twitch)
	if twitchUser == nil {
		return nil, fiber.NewError(http.StatusNotFound, "cannot find twitch user")
	}

	greeting.Enabled = *dto.Enabled
	err := services.DB.Model(greeting).Select("*").Updates(greeting).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot update greeting")
	}

	return &Greeting{
		ChannelsGreetings: *greeting,
		UserName:          twitchUser.Login,
		Avatar:            twitchUser.ProfileImageURL,
	}, nil
}
