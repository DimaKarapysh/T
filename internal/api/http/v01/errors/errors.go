package errors

import (
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Message string            `json:"message" example:"invalid request"`
	Reason  string            `json:"reason,omitempty" example:"invalid_guid"`
	Details map[string]string `json:"details,omitempty" example:"{\"field\":\"user_id\",\"description\":\"must be valid UUID\"}"`
}

func JSON(c *fiber.Ctx, status int, message string, reason string, details map[string]string) error {
	return c.Status(status).JSON(ErrorResponse{
		Message: message,
		Reason:  reason,
		Details: details,
	})
}

func ErrInvalidRequest(c *fiber.Ctx, reason string, details map[string]string) error {
	return JSON(c, fiber.StatusBadRequest, "invalid request", reason, details)
}

func ErrUnauthorized(c *fiber.Ctx, reason string) error {
	return JSON(c, fiber.StatusUnauthorized, "unauthorized", reason, nil)
}

func ErrInternal(c *fiber.Ctx, reason string, err error) error {
	details := make(map[string]string)
	if err != nil {
		details["error"] = err.Error()
	}
	return JSON(c, fiber.StatusInternalServerError, "internal error", reason, details)
}

func ErrNotFound(c *fiber.Ctx, reason string) error {
	return JSON(c, fiber.StatusNotFound, "not found", reason, nil)
}
