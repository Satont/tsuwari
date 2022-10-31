package keywords

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("keywords")

	middleware.Get("", get(services))
	middleware.Post("", post(services))
	middleware.Delete(":keywordId", delete(services))
	middleware.Patch(":keywordId", patch(services))

	return middleware
}

func get(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		keywords, err := handleGet(c.Params("channelId"), services)
		if err != nil {
			return err
		}
		return c.JSON(keywords)
	}
}

func post(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &keywordDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		keyword, err := handlePost(c.Params("channelId"), dto, services)
		if err == nil {
			return c.JSON(keyword)
		}

		return err
	}
}

func delete(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		err := handleDelete(c.Params("keywordId"), services)
		if err != nil {
			return err
		}

		return c.SendStatus(200)
	}
}

func patch(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &keywordDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		keyword, err := handleUpdate(c.Params("keywordId"), dto, services)
		if err != nil {
			return err
		}

		return c.JSON(keyword)
	}
}
