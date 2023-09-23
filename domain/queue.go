package domain

import (
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
	N1  float64
	D   float64
	I   float64
	TTL float64

	// Dynamic data with changing by workers
	Id               int
	CurrentIteration int
	Result           float64
	Status           TaskStatus
	StartedAt        time.Time
	EndedAt          time.Time
}

type IQueueService interface {
	AddTask(*Task) error
	GetTasks() []*Task
}
