package emotesbttv

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"
	"tsuwari/parser/pkg/helpers"

	"github.com/samber/lo"
)

type Emote struct {
	Code string `json:"code"`
}

type BttvResponse struct {
	ChannelEmotes []Emote `json:"channelEmotes"`
	SharedEmotes  []Emote `json:"sharedEmotes"`
}

var Variable = types.Variable{
	Name:        "emotes.bttv",
	Description: lo.ToPtr("Emotes of channel from https://betterttv.com/"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		resp, err := http.Get("https://api.betterttv.net/3/cached/users/twitch/" + ctx.ChannelId)
		if err != nil {
			log.Fatalln(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		reqData := BttvResponse{}
		err = json.Unmarshal(body, &reqData)
		if err != nil {
			return nil, errors.New("cannot fetch ffz emotes")
		}

		emotes := []string{}

		mappedChannelEmotes := helpers.Map(reqData.ChannelEmotes, func(e Emote) string {
			return e.Code
		})
		mappedSharedEmotes := helpers.Map(reqData.SharedEmotes, func(e Emote) string {
			return e.Code
		})

		emotes = append(emotes, mappedChannelEmotes...)
		emotes = append(emotes, mappedSharedEmotes...)

		result := types.VariableHandlerResult{
			Result: strings.Join(emotes, " "),
		}

		return &result, nil
	},
}
