package middleware

import (
	"github.com/BounkhongDev/bkgo/i18n"
	"github.com/BounkhongDev/bkgo/response"
	"github.com/gofiber/fiber/v2"
)

// RequireRole returns a middleware that allows only requests whose JWT claims
// contain a "role" value matching one of the provided roles.
// Must be used after the JWT middleware.
//
//	api.Post("/users", middleware.JWT(token), middleware.RequireRole("admin"), handler.Create)
func RequireRole(roles ...string) fiber.Handler {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}

	return func(c *fiber.Ctx) error {
		role, _ := Claims(c)["role"].(string)
		if !allowed[role] {
			locale := i18n.FromHeader(c.Get("Accept-Language"))
			return c.Status(fiber.StatusForbidden).JSON(
				response.Error("FORBIDDEN", i18n.Translate(locale, "FORBIDDEN")),
			)
		}
		return c.Next()
	}
}
