package domain

import (
	"context"
	"time"
)

type TaskStatus int

const (
	TaskPending    TaskStatus = 0
	TaskInProgress TaskStatus = 1
	TaskDone       TaskStatus = 2
)

type Task struct {
	// Static config data for runnable task
	N   int
	N1  int
	D   int
	I   int
	TTL int

	// Dynamic data with changing by workers
	Id               int
	CurrentIteration int
	Result           int
	Status           TaskStatus
	StartedAt        time.Time
	EndedAt          time.Time
}

type IQueueService interface {
	AddTask(context.Context, *Task) error
	GetTasks(context.Context) []*Task
}
