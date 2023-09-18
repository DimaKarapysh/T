package delivery

import (
	"T/model"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service model.IQueueService
}

func NewHandler(service model.IQueueService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Router(r *gin.RouterGroup) {
	r.POST("/reg", h.reg)
	r.GET("/get", h.get)
}

func (h *Handler) reg(r *gin.Context) {
	task := &model.Task{}
	err := r.ShouldBindJSON(task)
	if err != nil {
		_ = r.Error(NewUserError("Not invalid data", err))
		return
	}

	err = h.service.AddTask(task)
	if err != nil {
		_ = r.Error(err)
		r.JSON(200, err.Error())
		return
	}
	r.JSON(200, "success")
}

func (h *Handler) get(r *gin.Context) {
	r.JSON(200, h.service.GetJobs())
	return
}
