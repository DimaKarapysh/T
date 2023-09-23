package transport

import (
	"T/app/core"
	"T/app/rest"
	"T/domain"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"net/http"
)

type QueueTransportService struct {
	log     core.Logger
	service domain.IQueueService
}

func NewQueueTransportService(log core.Logger, s domain.IQueueService) *QueueTransportService {
	return &QueueTransportService{
		log:     log,
		service: s,
	}
}

func (t *QueueTransportService) AddTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		rest.ValidationError(w, "Cannot parse json")
		return
	}

	// ToDo: Make one point initialization (also need hide error message for prod)
	v := validator.New()
	err = v.Struct(task)
	if err != nil {
		rest.ValidationError(w, err.Error())
		return
	}

	err = t.service.AddTask(task.DTO())
	if err != nil {
		rest.ServerError(w, errors.Wrap(err, "AddTask"))
		return
	}

	rest.ServerSuccessOK(w)
}

func (t *QueueTransportService) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks := t.service.GetTasks()
	rest.ServerSuccessStruct(w, FromDomainTasks(tasks))
}
