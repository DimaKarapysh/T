package tests

import (
	"T/app"
	"T/domain"
	"T/services"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_PendingCheck(t *testing.T) {
	const N = 2

	ctx := context.Background()

	// init logs
	logger, err := app.InitLogs()
	require.NoError(t, err, "logger should be initialize")

	service := services.NewQueueService(logger, N)
	service.RunBackground()

	task1 := &domain.Task{
		N:   3,
		N1:  1,
		D:   1,
		I:   1,
		TTL: 5,
	}

	task2 := &domain.Task{
		N:   5,
		N1:  2,
		D:   1,
		I:   1,
		TTL: 0,
	}

	task3 := &domain.Task{
		N:   10,
		N1:  3,
		D:   1,
		I:   0,
		TTL: 3,
	}

	err = service.AddTask(ctx, task1)
	require.NoError(t, err, "task1 should be handled")

	err = service.AddTask(ctx, task2)
	require.NoError(t, err, "task2 should be handled")

	err = service.AddTask(ctx, task3)
	require.NoError(t, err, "task3 should be handled")

	// We have only N active processes (Checking that N+1 task is in pending status)
	time.Sleep(1 * time.Second)
	tasks := service.GetTasks(ctx)
	hasPending := false
	for _, task := range tasks {
		if task.Status == domain.TaskPending {
			hasPending = true
			break
		}
	}
	require.Equal(t, hasPending, true, "n+1 task should be in pending status")
}

func Test_CorrectAnswerCheck(t *testing.T) {
	ctx := context.Background()

	// init logs
	logger, err := app.InitLogs()
	require.NoError(t, err, "logger should be initialize")

	service := services.NewQueueService(logger, 1)
	service.RunBackground()

	task1 := &domain.Task{
		N:   3,
		N1:  5,
		D:   1,
		I:   1,
		TTL: 60,
	}

	err = service.AddTask(ctx, task1)
	require.NoError(t, err, "task1 should be handled")

	// Checking first tasks for correct answer
	time.Sleep(time.Duration(task1.N*task1.I+2) * time.Second)
	correct := task1.N1 + task1.D*task1.N

	tasks := service.GetTasks(ctx)
	require.Equal(t, len(tasks), 1, "tasks len should be == 1")
	require.Equal(t, tasks[0].Result, correct, "tasks answer is not correct")
}

func Test_TTLCheck(t *testing.T) {
	ctx := context.Background()

	// init logs
	logger, err := app.InitLogs()
	require.NoError(t, err, "logger should be initialize")

	service := services.NewQueueService(logger, 1)
	service.RunBackground()

	task1 := &domain.Task{
		N:   3,
		N1:  5,
		D:   1,
		I:   1,
		TTL: 5,
	}

	err = service.AddTask(ctx, task1)
	require.NoError(t, err, "task1 should be handled")

	// Checking first tasks for correct answer
	time.Sleep(time.Duration(task1.N*task1.I+2) * time.Second)
	tasks := service.GetTasks(ctx)
	require.Equal(t, len(tasks), 1, "tasks len should be == 1")

	time.Sleep(time.Duration(task1.TTL+5) * time.Second)
	tasks = service.GetTasks(ctx)
	require.Equal(t, len(tasks), 0, "tasks len should be 0")
}
