package v01

import (
	"T/internal/api/http/v01/errors"
	_ "T/internal/entity"
	"T/internal/service/adapters"
	"T/internal/validator"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Handler struct {
	service   adapters.SubService
	validator validator.Validator
	logger    *zap.Logger
}

func NewHandler(
	service adapters.SubService,
	validator validator.Validator,
	logger *zap.Logger,
) *Handler {
	return &Handler{
		service:   service,
		validator: validator,
		logger:    logger.With(zap.String("handler", "subscription")),
	}
}

// Create godoc
//
//	@Summary		Создать подписку
//	@Description	Создание новой подписки для пользователя
//	@Tags			Subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateSubscriptionRequest	true	"Тело запроса"
//	@Success		201		{object}	map[string]string	"ID созданной подписки"
//	@Failure		400		{object}	errors.ErrorResponse
//	@Failure		500		{object}	errors.ErrorResponse
//	@Router			/api/v1/subscriptions [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	const op = "Handler.Create"
	logger := h.logger.With(zap.String("op", op))

	var req CreateSubscriptionRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Warn("create: invalid body", zap.Error(err))
		return errors.ErrInvalidRequest(c, "invalid_body", nil)
	}

	if err := h.validator.Validate(req); err != nil {
		logger.Warn("create: validation failed", zap.Error(err))
		return errors.ErrInvalidRequest(c, "validation_failed", nil)
	}

	ctx := context.Background()

	toEntity, err := req.ToEntity()
	if err != nil {
		logger.Error("create: toEntity error", zap.Error(err))
		return errors.ErrInternal(c, "create_failed", err)
	}

	id, err := h.service.Create(ctx, toEntity)
	if err != nil {
		logger.Error("create: service error", zap.Error(err))
		return errors.ErrInternal(c, "create_failed", err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
}

// GetByID godoc
//
//	@Summary		Получить подписку по ID
//	@Description	Возвращает подписку по UUID
//	@Tags			Subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"ID подписки (UUID)"
//	@Success		200	{object}	entity.Subscription
//	@Failure		400	{object}	errors.ErrorResponse
//	@Failure		404	{object}	errors.ErrorResponse
//	@Failure		500	{object}	errors.ErrorResponse
//	@Router			/api/v1/subscriptions/{id} [get]
func (h *Handler) GetByID(c *fiber.Ctx) error {
	const op = "Handler.GetByID"
	logger := h.logger.With(zap.String("op", op))

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.ErrInvalidRequest(c, "invalid_id", nil)
	}

	ctx := context.Background()
	sub, err := h.service.GetByID(ctx, id)
	if err != nil {
		logger.Error("getById: not found", zap.Error(err))
		return errors.ErrNotFound(c, "not_found")
	}

	return c.JSON(sub)
}

// Update godoc
//
//	@Summary		Обновить подписку
//	@Description	Обновление полей существующей подписки
//	@Tags			Subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string						true	"ID подписки (UUID)"
//	@Param			request	body	UpdateSubscriptionRequest	true	"Тело запроса"
//	@Success		200		"OK"
//	@Failure		400		{object}	errors.ErrorResponse
//	@Failure		500		{object}	errors.ErrorResponse
//	@Router			/api/v1/subscriptions/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	const op = "Handler.Update"
	logger := h.logger.With(zap.String("op", op))

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Warn("update: invalid id", zap.Error(err))
		return errors.ErrInvalidRequest(c, "invalid_id", nil)
	}

	var req UpdateSubscriptionRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Warn("update: invalid body", zap.Error(err))
		return errors.ErrInvalidRequest(c, "invalid_body", nil)
	}

	if err := h.validator.Validate(req); err != nil {
		logger.Warn("update: validation failed", zap.Error(err))
		return errors.ErrInvalidRequest(c, "validation_failed", nil)
	}

	ctx := context.Background()

	toEntity, err := req.ToEntity()
	if err != nil {
		logger.Error("update: toEntity error", zap.Error(err))
		return errors.ErrInternal(c, "update_failed", err)
	}

	err = h.service.Update(ctx, id, toEntity)
	if err != nil {
		logger.Error("update: service error", zap.Error(err))
		return errors.ErrInternal(c, "update_failed", err)
	}

	return c.SendStatus(fiber.StatusOK)
}

// Delete godoc
//
//	@Summary		Удалить подписку
//	@Description	Удаляет подписку по её UUID (мягкое или жёсткое удаление зависит от реализации)
//	@Tags			Subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID подписки (UUID)"
//	@Success		200	"OK"
//	@Failure		400	{object}	errors.ErrorResponse
//	@Failure		500	{object}	errors.ErrorResponse
//	@Router			/api/v1/subscriptions/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	const op = "Handler.Delete"
	logger := h.logger.With(zap.String("op", op))

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Warn("delete: invalid id", zap.Error(err))
		return errors.ErrInvalidRequest(c, "invalid_id", nil)
	}

	ctx := context.Background()
	if err := h.service.Delete(ctx, id); err != nil {
		logger.Error("delete: service error", zap.Error(err))
		return errors.ErrInternal(c, "delete_failed", err)
	}

	return c.SendStatus(fiber.StatusOK)
}

// List godoc
//
//	@Summary		Получение списка подписок
//	@Description	Пагинированный список подписок
//	@Tags			Subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	false	"Максимум записей"		minimum(1) maximum(100) default(20)
//	@Param			offset	query		int	false	"Смещение"				minimum(0) default(0)
//	@Success		200		{array}		entity.Subscription
//	@Failure		400		{object}	errors.ErrorResponse
//	@Failure		500		{object}	errors.ErrorResponse
//	@Router			/api/v1/subscriptions [get]
func (h *Handler) List(c *fiber.Ctx) error {
	const op = "Handler.List"

	logger := h.logger.With(zap.String("op", op))

	var query ListQueryParams
	if err := c.QueryParser(&query); err != nil {
		logger.Error("ошибка парсинга query", zap.Error(err))
		return errors.ErrInvalidRequest(c, "invalid_query", nil)
	}

	// дефолты
	limit := 20
	if query.Limit != nil {
		limit = *query.Limit
		if limit <= 0 || limit > 100 {
			logger.Error("limit must be between 1 and 100")
			return errors.ErrInvalidRequest(c, "limit must be between 1 and 100", nil)
		}
	}

	offset := 0
	if query.Offset != nil {
		offset = *query.Offset
		if offset < 0 {
			logger.Error("offset must be >= 0")
			return errors.ErrInvalidRequest(c, "offset must be >= 0", nil)
		}
	}

	logger.Info("получение подписок", zap.Int("limit", limit), zap.Int("offset", offset))

	subs, err := h.service.List(c.Context(), limit, offset)
	if err != nil {
		logger.Error("ошибка получения списка", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	return c.JSON(subs)
}

// Total godoc
//
//	@Summary		Получение общей стоимости подписок
//	@Description	Фильтрация по пользователю, сервису и периоду (через JSON)
//	@Tags			Subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			request	body		TotalRequest	true	"Фильтры"
//	@Success		200		{object}	map[string]int
//	@Failure		400		{object}	errors.ErrorResponse
//	@Failure		500		{object}	errors.ErrorResponse
//	@Router			/api/v1/subscriptions/total [post]
func (h *Handler) Total(c *fiber.Ctx) error {
	const op = "Handler.Total"

	logger := h.logger.With(zap.String("op", op))

	var req TotalRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Warn("total: invalid body", zap.Error(err))
		return errors.ErrInvalidRequest(c, "invalid_body", nil)
	}

	if err := h.validator.Validate(req); err != nil {
		return errors.ErrInvalidRequest(c, "validation_failed", nil)
	}

	filter, err := req.ToFilter()
	if err != nil {
		logger.Warn("total: invalid input", zap.Error(err))
		return errors.ErrInvalidRequest(c, "invalid_input", nil)
	}

	ctx := context.Background()
	total, err := h.service.Total(ctx, filter)
	if err != nil {
		logger.Error("total: service error", zap.Error(err))
		return errors.ErrInternal(c, "total_failed", err)
	}

	return c.JSON(fiber.Map{"total_price": total})
}
