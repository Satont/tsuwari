package manage

import (
	"context"
	"fmt"
	"log"
	"strings"
	model "tsuwari/parser/internal/models"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/pkg/helpers"

	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/guregu/null"
	"github.com/samber/lo"
	uuid "github.com/satori/go.uuid"
)

const (
	exampleUsage = "!commands add name response"
	incorrectUsage = "Incorrect usage of command. Example: " + exampleUsage
	wentWrong = "Something went wrong on creating command"
	alreadyExists = "Command with that name or aliase already exists."
)

var AddCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "commands add",
		Description: lo.ToPtr("Add command"),
		Permission:  "MODERATOR",
		Visible:     true,
		Module:      lo.ToPtr("MANAGE"),
	},
	Handler: func(ctx variables_cache.ExecutionContext) []string {
		if ctx.Text == nil {
			return []string{incorrectUsage}
		}

		args := strings.Split(*ctx.Text, " ")

		if len(args) < 2 {
			return []string{incorrectUsage}
		}

		name := args[0]
		text := strings.Join(args[1:], " ")

		if len(name) > 20 {
			return []string{"Command name cannot be greatest then 20."}
		}

		commands := []model.ChannelsCommands{}
		err := ctx.Services.Db.Model(&model.ChannelsCommands{}).
			Where(`"channelId" = ?`, ctx.ChannelId).
			Find(&commands).Error

		if err != nil {
			log.Fatalln(err)
			return []string{}
		}

		for _, c := range commands {
			if c.Name == name {
				return []string{alreadyExists}
			}

			if helpers.Contains(c.Aliases, name) {
				return []string{alreadyExists}
			}
		}

		commandID := uuid.NewV4().String()
		command := model.ChannelsCommands{
			ID: commandID,
			Name: name,
			CooldownType: "GLOBAL",
			Enabled: true,
			Cooldown: null.IntFrom(5),
			Aliases: []string{},
			Description: null.String{},
			DefaultName: null.String{},
			Visible: true,
			ChannelID: ctx.ChannelId,
			Permission: "VIEWER",
			Default: false,
			Module: "CUSTOM",
			Responses: []*model.ChannelsCommandsResponses{
				{
					ID: uuid.NewV4().String(),
					Text: null.StringFrom(text),
					CommandID: commandID,
				},
			},
		}
		err = ctx.Services.Db.Create(&command).Error

		if err != nil {
			log.Fatalln(err)
			return []string{wentWrong}
		}


		bytes, err := CreateRedisBytes(command, text, lo.ToPtr(true))

		_, err = ctx.Services.Redis.Set(
			context.TODO(), 
			fmt.Sprintf("commands:%s:%s", ctx.ChannelId, name),
			*bytes,
			0,
		).Result()

		if err != nil {
			log.Fatalln("cannot create command in redis", err)
			return []string{wentWrong}
		}

		ctx.Services.Redis.Del(
			context.TODO(), 
			fmt.Sprintf("nest:cache:v1/channels/%s/commands", ctx.ChannelId), 
		)


		return []string{"✅ Command added."}
	},
}