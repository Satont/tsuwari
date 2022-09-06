package faceitelodiff

import (
	"strconv"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "faceit.todayEloDiff",
	Description: lo.ToPtr("Faceit today elo earned"),
	Handler: func(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		matches := ctx.GetFaceitLatestMatches()
		diff := ctx.GetFaceitTodayEloDiff(matches)

		result.Result = strconv.Itoa(diff)

		return result, nil
	},
}
