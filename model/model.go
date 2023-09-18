package model

import "time"

type Task struct {
	N   int     `json:"n"`
	D   float64 `json:"d"`
	N1  float64 `json:"n_1"`
	I   float64 `json:"i"`
	TTL float64 `json:"ttl"`
}

type Job struct {
	Task          *Task     `json:"task"`
	State         string    `json:"state"`
	NumberOfQueue int       `json:"number_of_queue"`
	CurrentValue  float64   `json:"current_value"`
	CurrentI      int       `json:"itteration_n"`
	StandedAt     time.Time `json:"standedAt"`
	StartedAt     time.Time `json:"startedAt"`
	EndedAt       time.Time `json:"endedAt"`
}

type IQueueService interface {
	AddTask(d *Task) error
	GetJobs() []*Job
}
