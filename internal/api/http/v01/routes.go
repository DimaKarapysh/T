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
	Total(c *fiber.Ctx) error // Сумма подписок по фильтрам
}

func RegisterRoutes(router fiber.Router, h HttpHandler) {
	router.Post("/subscriptions", h.Create)
	router.Get("/subscriptions/:id", h.GetByID)
	router.Put("/subscriptions/:id", h.Update)
	router.Delete("/subscriptions/:id", h.Delete)
	router.Get("/subscriptions", h.List)

	// ручка для подсчета общей стоимости по фильтрам
	router.Post("/subscriptions/total", h.Total)
}
