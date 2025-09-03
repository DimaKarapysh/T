package v01

import (
	"T/internal/api/http/v01/errors"
	_ "T/internal/entity"
	"T/internal/service/adapters"
	"T/internal/util"
	"T/internal/validator"
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Handler struct {
	service   adapters.SubService
	validator validator.Validator
	logger    *zap.Logger
	tr        trace.Tracer
}

func NewHandler(
	service adapters.SubService,
	validator validator.Validator,
	logger *zap.Logger,
	tracer trace.Tracer,
) *Handler {
	return &Handler{
		service:   service,
		validator: validator,
		logger:    logger.With(zap.String("handler", "subscription")),
		tr:        tracer,
	}
}

// Get godoc
//
//	@Summary		Получить все задачи
//	@Description	Возвращает список всех задач без пагинации
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		entity.Task
//	@Failure		500	{object}	errors.ErrorResponse
//	@Router			/api/v1/tasks [get]
func (h *Handler) Get(c *fiber.Ctx) error {
	const op = "Handler.Get"
	logger := h.logger.With(zap.String("op", op))

	ctx := context.Background()
	ctx, step := util.Start(ctx, h.tr, logger, op)
	var spanErr error
	defer step.End(spanErr)

	tasks, err := h.service.Get(ctx)
	if err != nil {
		logger.Error("ошибка получения всех задач", zap.Error(err))
		spanErr = err
		return errors.ErrInternal(c, "get_failed", err)
	}

	return c.JSON(tasks)
}

// Create godoc
//
//	@Summary		Создать задачу
//	@Description	Создание новой задачи для пользователя
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateTaskRequest	true	"Тело запроса"
//	@Success		201		{object}	map[string]string	"ID созданной задачи"
//	@Failure		400		{object}	errors.ErrorResponse
//	@Failure		500		{object}	errors.ErrorResponse
//	@Router			/api/v1/tasks [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	const op = "Handler.Create"
	logger := h.logger.With(zap.String("op", op))

	ctx := context.Background()
	ctx, step := util.Start(ctx, h.tr, logger, op)
	var spanErr error
	defer step.End(spanErr)

	var req CreateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Warn("create: invalid body", zap.Error(err))
		spanErr = err
		return errors.ErrInvalidRequest(c, "invalid_body", nil)
	}

	if err := h.validator.Validate(req); err != nil {
		logger.Warn("create: validation failed", zap.Error(err))
		spanErr = err
		return errors.ErrInvalidRequest(c, "validation_failed", nil)
	}

	toEntity, err := req.ToEntity()
	if err != nil {
		logger.Error("create: toEntity error", zap.Error(err))
		spanErr = err
		return errors.ErrInternal(c, "create_failed", err)
	}

	id, err := h.service.Create(ctx, toEntity)
	if err != nil {
		logger.Error("create: service error", zap.Error(err))
		spanErr = err
		return errors.ErrInternal(c, "create_failed", err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
}

// GetByID godoc
//
//	@Summary		Получить задачу по ID
//	@Description	Возвращает задачу по UUID
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"ID задачи (UUID)"
//	@Success		200	{object}	entity.Task
//	@Failure		400	{object}	errors.ErrorResponse
//	@Failure		404	{object}	errors.ErrorResponse
//	@Failure		500	{object}	errors.ErrorResponse
//	@Router			/api/v1/tasks/{id} [get]
func (h *Handler) GetByID(c *fiber.Ctx) error {
	const op = "Handler.GetByID"
	logger := h.logger.With(zap.String("op", op))

	ctx := context.Background()
	ctx, step := util.Start(ctx, h.tr, logger, op)
	var spanErr error
	defer step.End(spanErr)

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		spanErr = err
		return errors.ErrInvalidRequest(c, "invalid_id", nil)
	}

	task, err := h.service.GetByID(ctx, id)
	if err != nil {
		logger.Error("getById: not found", zap.Error(err))
		spanErr = err
		return errors.ErrNotFound(c, "not_found")
	}

	return c.JSON(task)
}

// Update godoc
//
//	@Summary		Обновить задачу
//	@Description	Обновление полей существующей задачи
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string				true	"ID задачи (UUID)"
//	@Param			request	body	UpdateTaskRequest	true	"Тело запроса"
//	@Success		200		"OK"
//	@Failure		400		{object}	errors.ErrorResponse
//	@Failure		500		{object}	errors.ErrorResponse
//	@Router			/api/v1/tasks/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	const op = "Handler.Update"
	logger := h.logger.With(zap.String("op", op))

	ctx := context.Background()
	ctx, step := util.Start(ctx, h.tr, logger, op)
	var spanErr error
	defer step.End(spanErr)

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Warn("update: invalid id", zap.Error(err))
		spanErr = err
		return errors.ErrInvalidRequest(c, "invalid_id", nil)
	}

	var req UpdateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Warn("update: invalid body", zap.Error(err))
		spanErr = err
		return errors.ErrInvalidRequest(c, "invalid_body", nil)
	}
	if err := h.validator.Validate(req); err != nil {
		logger.Warn("update: validation failed", zap.Error(err))
		spanErr = err
		return errors.ErrInvalidRequest(c, "validation_failed", nil)
	}

	toEntity, err := req.ToEntity()
	if err != nil {
		logger.Error("update: toEntity error", zap.Error(err))
		spanErr = err
		return errors.ErrInternal(c, "update_failed", err)
	}

	if err := h.service.Update(ctx, id, toEntity); err != nil {
		logger.Error("update: service error", zap.Error(err))
		spanErr = err
		return errors.ErrInternal(c, "update_failed", err)
	}

	return c.SendStatus(fiber.StatusOK)
}

// Delete godoc
//
//	@Summary		Удалить задачу
//	@Description	Удаляет задачу по её UUID (мягкое или жёсткое удаление зависит от реализации)
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID задачи (UUID)"
//	@Success		200	"OK"
//	@Failure		400	{object}	errors.ErrorResponse
//	@Failure		500	{object}	errors.ErrorResponse
//	@Router			/api/v1/tasks/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	const op = "Handler.Delete"
	logger := h.logger.With(zap.String("op", op))

	ctx := context.Background()
	ctx, step := util.Start(ctx, h.tr, logger, op)
	var spanErr error
	defer step.End(spanErr)

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Warn("delete: invalid id", zap.Error(err))
		spanErr = err
		return errors.ErrInvalidRequest(c, "invalid_id", nil)
	}

	if err := h.service.Delete(ctx, id); err != nil {
		logger.Error("delete: service error", zap.Error(err))
		spanErr = err
		return errors.ErrInternal(c, "delete_failed", err)
	}

	logger.Info("delete: ok", zap.String("id", id.String()))
	return c.SendStatus(fiber.StatusOK)
}

// List godoc
//
//	@Summary		Получение списка задач
//	@Description	Пагинированный список задач
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	false	"Максимум записей"		minimum(1) maximum(100) default(20)
//	@Param			offset	query		int	false	"Смещение"				minimum(0) default(0)
//	@Success		200		{array}		entity.Task
//	@Failure		400		{object}	errors.ErrorResponse
//	@Failure		500		{object}	errors.ErrorResponse
//	@Router			/api/v1/tasks [get]
func (h *Handler) List(c *fiber.Ctx) error {
	const op = "Handler.List"
	logger := h.logger.With(zap.String("op", op))

	ctx := context.Background()
	ctx, step := util.Start(ctx, h.tr, logger, op)
	var spanErr error
	defer step.End(spanErr)

	var query ListQueryParams
	if err := c.QueryParser(&query); err != nil {
		logger.Error("ошибка парсинга query", zap.Error(err))
		spanErr = err
		return errors.ErrInvalidRequest(c, "invalid_query", nil)
	}

	limit := 20
	if query.Limit != nil {
		limit = *query.Limit
		if limit <= 0 || limit > 100 {
			err := fmt.Errorf("limit must be between 1 and 100")
			logger.Error(err.Error())
			spanErr = err
			return errors.ErrInvalidRequest(c, err.Error(), nil)
		}
	}
	offset := 0
	if query.Offset != nil {
		offset = *query.Offset
		if offset < 0 {
			err := fmt.Errorf("offset must be >= 0")
			logger.Error(err.Error())
			spanErr = err
			return errors.ErrInvalidRequest(c, err.Error(), nil)
		}
	}

	tasks, err := h.service.List(ctx, limit, offset)
	if err != nil {
		logger.Error("ошибка получения списка", zap.Error(err))
		spanErr = err
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}
	return c.JSON(tasks)
}
