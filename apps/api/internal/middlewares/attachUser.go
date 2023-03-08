package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

var AttachUser = func(services *types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Locals("dbUser") != nil {
			return c.Next()
		}

		headers := c.GetReqHeaders()
		apiKey := headers["Api-Key"]
		dbUser := model.Users{}

		if apiKey != "" {
			err := services.Gorm.
				Where(`"apiKey" = ?`, apiKey).
				Preload("Roles").
				Preload("Roles.Role").
				First(&dbUser).Error
			if err != nil {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "user with that api key not found"})
			}
			c.Locals("dbUser", dbUser)
			return c.Next()
		}

		authorizationToken := headers["Authorization"]
		if authorizationToken != "" {
			token, err := ExtractTokenFromHeader(services, authorizationToken)
			if err != nil {
				return fiber.NewError(http.StatusUnauthorized, "invalid token. Probably token is expired.")
			}

			claims := token.Claims.(jwt.MapClaims)
			userId := claims["id"]

			if userId == "" {
				services.Logger.Error("no userId in request")
				return fiber.NewError(http.StatusUnauthorized, "invalid token")
			}

			err = services.Gorm.
				Where(`"id" = ?`, userId).
				Preload("Roles").
				Preload("Roles.Role").
				Find(&dbUser).Error
			if err != nil {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "user not found"})
			}
			c.Locals("dbUser", dbUser)
		}

		if dbUser.ID == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "unauthenticated"})
		}

		return c.Next()
	}
}

func ExtractTokenFromHeader(services *types.Services, t string) (*jwt.Token, error) {
	tokenSlice := strings.Split(t, "Bearer ")
	if len(tokenSlice) < 2 {
		return nil, fiber.NewError(http.StatusUnauthorized, "invalid token format")
	}

	token, err := jwt.Parse(tokenSlice[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(services.Config.JwtAccessSecret), nil
	})

	return token, err
}
