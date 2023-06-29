package valorant

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/twir/apps/api/internal/middlewares"
	"github.com/satont/twir/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("valorant")
	middleware.Get("", get(services))
	middleware.Post("", post(services))

	return middleware
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Get valorant data
// @Tags         Integrations|Valorant
// @Accept       json
// @Produce      plain
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200  {object} model.ChannelsIntegrationsData
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/valorant [get]
func get(services types.Services) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		data, err := handleGet(services, ctx.Params("channelId"))

		if err != nil {
			return err
		}

		return ctx.JSON(data)
	}
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Update valorant data
// @Tags         Integrations|Valorant
// @Accept       json
// @Param data body createOrUpdateDTO true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/valorant [post]
func post(services types.Services) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		dto := &createOrUpdateDTO{}
		err := middlewares.ValidateBody(
			ctx,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		err = handlePost(services, ctx.Params("channelId"), dto)
		if err != nil {
			return err
		}

		return ctx.SendStatus(200)
	}
}
