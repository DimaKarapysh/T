package middlewarepackage

import "github.com/gofiber/fiber/v2"

type JWTMiddleware interface {
	AuthMiddleware() fiber.Handler
}
