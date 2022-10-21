package apiv1

import (
	"github.com/gofiber/fiber/v2"

	"github.com/satont/tsuwari/apps/api-go/internal/api/v1/commands"
	"github.com/satont/tsuwari/apps/api-go/internal/api/v1/greetings"
	"github.com/satont/tsuwari/apps/api-go/internal/api/v1/keywords"
	"github.com/satont/tsuwari/apps/api-go/internal/api/v1/moderation"
	"github.com/satont/tsuwari/apps/api-go/internal/api/v1/timers"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	channelsGroup := router.Group("channels/:channelId")
	commands.Setup(channelsGroup, services)
	greetings.Setup(channelsGroup, services)
	keywords.Setup(channelsGroup, services)
	timers.Setup(channelsGroup, services)
	moderation.Setup(channelsGroup, services)

	return router
}
