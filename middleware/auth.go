package middleware

import (
	"strings"

	"github.com/bounkhongdev/kbgo/contract"
	"github.com/bounkhongdev/kbgo/i18n"
	"github.com/bounkhongdev/kbgo/response"
	"github.com/gofiber/fiber/v2"
)

const claimsKey = "claims"

// JWT returns a Fiber middleware that validates a Bearer token from the
// Authorization header using the provided contract.Token adapter.
// On success it stores the claims in Fiber locals — retrieve with Claims(c).
func JWT(token contract.Token) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return unauthorized(c)
		}

		claims, err := token.Verify(strings.TrimPrefix(auth, "Bearer "))
		if err != nil {
			return unauthorized(c)
		}

		c.Locals(claimsKey, claims)
		return c.Next()
	}
}

// Claims retrieves the JWT claims stored by the JWT middleware.
// Returns an empty map if the middleware has not run or the token was invalid.
func Claims(c *fiber.Ctx) contract.Claims {
	if claims, ok := c.Locals(claimsKey).(contract.Claims); ok {
		return claims
	}
	return contract.Claims{}
}

func unauthorized(c *fiber.Ctx) error {
	locale := i18n.FromHeader(c.Get("Accept-Language"))
	return c.Status(fiber.StatusUnauthorized).JSON(
		response.Error("UNAUTHORIZED", i18n.Translate(locale, "UNAUTHORIZED")),
	)
}
