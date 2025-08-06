package http

import (
	v01 "T/internal/api/http/v01"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App, handler *v01.Handler) {
	// Swagger UI
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/docs/swagger.yaml",
	}))
	app.Static("/docs", "./docs")

	api := app.Group("/api/v1")

	v01.RegisterRoutes(api, handler)

	api.Get("/health", func(c *fiber.Ctx) error { // <--- вот тут
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "auth_service",
			"version": "1.0.0",
		})
	})

	// 404 fallback
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"errwrap": "endpoint not found",
		})
	})

	for _, route := range app.GetRoutes() {
		fmt.Printf("Route: %s %s\n", route.Method, route.Path)
	}
}
