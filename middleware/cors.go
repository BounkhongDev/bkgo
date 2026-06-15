package middleware

import "github.com/gofiber/fiber/v2"

// CORSConfig holds CORS policy settings.
type CORSConfig struct {
	AllowOrigins string // comma-separated origins, "*" for all
	AllowMethods string
	AllowHeaders string
}

var defaultCORSConfig = CORSConfig{
	AllowOrigins: "*",
	AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
	AllowHeaders: "Origin,Content-Type,Authorization,Accept-Language",
}

// CORS returns a Fiber middleware that sets CORS response headers.
// Call with no arguments to use permissive defaults, or pass a CORSConfig.
//
//	app.Use(middleware.CORS())
//	app.Use(middleware.CORS(middleware.CORSConfig{AllowOrigins: "https://myapp.com"}))
func CORS(cfg ...CORSConfig) fiber.Handler {
	c := defaultCORSConfig
	if len(cfg) > 0 {
		c = cfg[0]
	}

	return func(ctx *fiber.Ctx) error {
		ctx.Set("Access-Control-Allow-Origin", c.AllowOrigins)
		ctx.Set("Access-Control-Allow-Methods", c.AllowMethods)
		ctx.Set("Access-Control-Allow-Headers", c.AllowHeaders)

		if ctx.Method() == fiber.MethodOptions {
			return ctx.SendStatus(fiber.StatusNoContent)
		}
		return ctx.Next()
	}
}
