package commands

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("commands")

	middleware.Get("", get(services))
	middleware.Post("", post(services))
	middleware.Delete(":commandId", delete(services))
	middleware.Put(":commandId", put(services))

	return router
}

type JSONResult struct{}

// Commands godoc
// @Security ApiKeyAuth
// @Summary      Get channel commands list
// @Tags         Commands
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId"
// @Success      200  {array}  model.ChannelsCommands
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/commands [get]
func get(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.JSON(handleGet(c.Params("channelId"), services))

		return nil
	}
}

// Commands godoc
// @Security ApiKeyAuth
// @Summary      Create command
// @Tags         Commands
// @Accept       json
// @Produce      json
// @Param data body commandDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200  {object}  model.ChannelsCommands
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/commands [post]
func post(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &commandDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		cmd, err := handlePost(c.Params("channelId"), services, dto)
		if err == nil {
			return c.JSON(cmd)
		}

		return err
	}
}

// Commands godoc
// @Security ApiKeyAuth
// @Summary      Delete command
// @Tags         Commands
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ID of channel"
// @Param        commandId   path      string  true  "ID of command"
// @Success      200  {object}  model.ChannelsCommands
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/commands/{commandId} [delete]
func delete(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		err := handleDelete(c.Params("channelId"), c.Params("commandId"), services)
		if err != nil {
			return err
		}
		return c.SendStatus(fiber.StatusOK)
	}
}

// Commands godoc
// @Security ApiKeyAuth
// @Summary      Update command
// @Tags         Commands
// @Accept       json
// @Produce      json
// @Param data body commandDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Param        commandId   path      string  true  "ID of command"
// @Success      200  {object}  model.ChannelsCommands
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/commands/{commandId} [put]
func put(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &commandDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		cmd, err := handleUpdate(c.Params("channelId"), c.Params("commandId"), dto, services)
		if err == nil && cmd != nil {
			return c.JSON(cmd)
		}

		return err
	}
}
