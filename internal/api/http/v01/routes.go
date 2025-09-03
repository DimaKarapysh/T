package v01

import "github.com/gofiber/fiber/v2"

type HttpHandler interface {
	SubscriptionHandler
}

type SubscriptionHandler interface {
	Create(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	List(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
}

func RegisterRoutes(router fiber.Router, h HttpHandler) {
	router.Post("/tasks", h.Create)
	router.Get("/tasks/:id", h.GetByID)
	router.Put("/tasks/:id", h.Update)
	router.Delete("/tasks/:id", h.Delete)
	router.Get("/tasks", h.List)
	router.Get("/task", h.Get)

}
