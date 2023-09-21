package transport

import (
	"T/domain"
	"time"
)

type Task struct {
	N   int `json:"n" validate:"required"`
	N1  int `json:"n1"`
	D   int `json:"d"`
	I   int `json:"i"`
	TTL int `json:"ttl"`
}

func (t *Task) DTO() *domain.Task {
	return &domain.Task{
		N:   t.N,
		N1:  t.N1,
		D:   t.D,
		I:   t.I,
		TTL: t.TTL,
	}
}

type TaskInfo struct {
	Id               int       `json:"id"`
	N                int       `json:"n"`
	N1               int       `json:"n1"`
	D                int       `json:"d"`
	I                int       `json:"i"`
	TTL              int       `json:"ttl"`
	CurrentIteration int       `json:"current_iteration"`
	Result           int       `json:"current_result"`
	Status           string    `json:"status"`
	StartedAt        time.Time `json:"started_at"`
	EndedAt          time.Time `json:"ended_at"`
}

func FromDomainTask(task *domain.Task) *TaskInfo {

	status := ""
	switch task.Status {
	case domain.TaskPending:
		status = "В ожидании"
	case domain.TaskInProgress:
		status = "В процессе"
	default:
		status = "Выполнена"
	}

	return &TaskInfo{
		Id:               task.Id,
		N:                task.N,
		N1:               task.N1,
		D:                task.D,
		I:                task.I,
		TTL:              task.TTL,
		CurrentIteration: task.CurrentIteration,
		Result:           task.Result,
		Status:           status,
		StartedAt:        task.StartedAt,
		EndedAt:          task.EndedAt,
	}
}

func FromDomainTasks(tasks []*domain.Task) []*TaskInfo {
	var taskInfos []*TaskInfo
	for _, t := range tasks {
		taskInfos = append(taskInfos, FromDomainTask(t))
	}
	return taskInfos
}
