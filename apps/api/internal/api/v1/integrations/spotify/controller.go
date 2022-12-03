package spotify

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/satont/tsuwari/apps/api/internal/types"

	"github.com/satont/tsuwari/apps/api/internal/middlewares"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("spotify")
	middleware.Get("auth", getAuth(services))
	middleware.Get("", get(services))
	middleware.Post("token", post((services)))
	middleware.Patch("", patch((services)))

	profileCache := cache.New(cache.Config{
		Expiration: 31 * 24 * time.Hour,
		Storage:    services.RedisStorage,
		KeyGenerator: func(c *fiber.Ctx) string {
			return fmt.Sprintf("fiber:cache:integrations:spotify:profile:%s", c.Params("channelId"))
		},
	})

	middleware.Get("profile", profileCache, getProfile((services)))

	return middleware
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Get Spotify integration
// @Tags         Integrations|Spotify
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId"
// @Success      200  {object}  model.ChannelsIntegrations
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/spotify [get]
func get(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		integration, err := handleGet(c.Params("channelId"), services)
		if err != nil {
			return err
		}
		return c.JSON(integration)
	}
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Get Spotify profile
// @Tags         Integrations|Spotify
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId"
// @Success      200  {object}  spotify.SpotifyProfile
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/spotify/profile [get]
func getProfile(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		profile, err := handleGetProfile(c.Params("channelId"), services)
		if err != nil {
			return err
		}
		return c.JSON(profile)
	}
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Get Spotify auth link
// @Tags         Integrations|Spotify
// @Accept       json
// @Produce      plain
// @Param        channelId   path      string  true  "ChannelId"
// @Success 200 {string} string	"Auth link"
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/spotify/auth [get]
func getAuth(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authLink, err := handleGetAuth(services)
		if err != nil {
			return err
		}

		return c.SendString(*authLink)
	}
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Update Spotify status
// @Tags         Integrations|Spotify
// @Accept       json
// @Produce      json
// @Param data body spotifyDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200  {object} model.ChannelsIntegrations
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/spotify [patch]
func patch(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &spotifyDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		integration, err := handlePatch(c.Params("channelId"), dto, services)
		if err != nil {
			return err
		}

		return c.JSON(integration)
	}
}

// Integrations godoc
// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Update auth of Spotify
// @Tags         Integrations|Spotify
// @Accept       json
// @Produce      json
// @Param data body tokenDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/spotify/token [post]
func post(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &tokenDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		channelId := c.Params("channelId")
		err = handlePost(channelId, dto, services)
		if err != nil {
			return err
		}

		services.RedisStorage.DeleteByMethod(
			fmt.Sprintf("fiber:cache:integrations:spotify:profile:%s_GET", channelId),
			"GET",
		)

		return c.SendStatus(200)
	}
}
