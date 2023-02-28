package feedback

import (
	"fmt"
	"time"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("feedback")
	middleware.Use(middlewares.AttachUser(services))

	limit := limiter.New(limiter.Config{
		Max:        2,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			dbUser := c.Locals("dbUser").(model.Users)
			return fmt.Sprintf("fiber:limiter:feedback:%s", dbUser.ID)
		},
		LimitReached: func(c *fiber.Ctx) error {
			header := c.GetRespHeader("Retry-After", "2")
			return c.Status(429).JSON(fiber.Map{"message": fmt.Sprintf("wait %s seconds", header)})
		},
		Storage: services.RedisStorage,
	})

	middleware.Post("", limit, func(c *fiber.Ctx) error {
		form, err := c.MultipartForm()
		if err != nil {
			return fiber.NewError(400, "wrong form data")
		}

		if form.Value["text"] == nil {
			return fiber.NewError(400, "text cannot be empty")
		}

		// TODO: add user id to the params
		err = handlePost("", form.Value["text"][0], form.File["files"], services)
		if err != nil {
			return err
		}

		return c.SendStatus(200)
	})

	return middleware
}
