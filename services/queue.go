package services

import (
	"T/app/core"
	"T/domain"
	"context"
	"sync"
	"time"
)

type QueueService struct {
	log    core.Logger
	N      int
	iterId int

	// Channels for communication
	tasksIn    chan *domain.Task
	tasksInRec chan *domain.Task

	// Current tasks with all data
	mu           sync.Mutex
	tasksRecords []*domain.Task
}

func NewQueueService(log core.Logger, n int) *QueueService {
	return &QueueService{
		log:          log,
		N:            n,
		tasksIn:      make(chan *domain.Task, n),
		tasksInRec:   make(chan *domain.Task, 1000), // 1000 - optional
		tasksRecords: []*domain.Task{},
	}
}

// core logic for tasks handling
func (q *QueueService) worker() {
	for task := range q.tasksIn {
		task.Status = domain.TaskInProgress
		task.StartedAt = time.Now()
		task.CurrentIteration = 0
		task.Result = task.N1
		q.log.Info("Task (%d) in PROGRESS", task.Id)

		for {
			// End of task
			if task.CurrentIteration >= task.N {
				task.Status = domain.TaskDone
				task.EndedAt = time.Now()
				q.log.Info("Task (%d) DONE", task.Id)
				break
			}

			task.Result += task.D
			task.CurrentIteration++
			q.log.Info("Task (%d) iteration %d. Value=%v", task.Id, task.CurrentIteration, task.Result)

			// waiting I seconds before next iteration
			time.Sleep(time.Duration(task.I) * time.Second)
		}
	}
}

func (q *QueueService) RunBackground() {
	for w := 1; w <= q.N; w++ {
		go q.worker()
	}

	go func() {
		for {
			for task := range q.tasksInRec {
				q.mu.Lock()
				q.tasksRecords = append(q.tasksRecords, task)
				q.mu.Unlock()
			}
		}
	}()

	// TTL expired
	go func() {
		// Checking for TTL
		for {
			q.mu.Lock()
			var valid []*domain.Task
			for _, task := range q.tasksRecords {
				if task.Status == domain.TaskDone &&
					time.Now().Add(-5*time.Second).After(task.EndedAt.Add(time.Second*time.Duration(task.TTL))) {
					q.log.Info("Task (%d) TTL expired", task.Id)
				} else {
					valid = append(valid, task)
				}
			}
			q.tasksRecords = valid
			q.mu.Unlock()

			time.Sleep(2 * time.Second)
		}
	}()
}

func (q *QueueService) AddTask(ctx context.Context, task *domain.Task) error {
	q.iterId++
	task.Status = domain.TaskPending
	task.Id = q.iterId

	q.tasksInRec <- task

	q.log.Info("Task (%d) in PENDING", task.Id)
	q.tasksIn <- task
	return nil
}

func (q *QueueService) GetTasks(ctx context.Context) []*domain.Task {
	return q.tasksRecords
}
