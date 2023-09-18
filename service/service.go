package service

import (
	"T/model"
	"log"
	"strconv"
	"strings"
	"time"
)

var count int

func process(s *QueueService) {
	for {
		// Если есть свободный ресурс, найти ожидающие запуск задачи и запустить
		if count < s.N {

			var job *model.Job
			for i := range s.Queue {
				if s.Queue[i].State == "In Queue" {
					job = s.Queue[i]
					break
				}
			}

			// Если ресурс свободен и найдена задача, требующая запуска
			if job != nil {
				count++ // Увеличивает счетчик активных процессов
				log.Println("New task in process")
				log.Println("Queue len=" + strconv.Itoa(len(s.Queue)))
				go func(job *model.Job) {
					for {
						if job.CurrentI >= job.Task.N {
							job.EndedAt = time.Now()
							job.State = "end"
							count--
							log.Println("Task done. Queue len=" + strconv.Itoa(len(s.Queue)))
							return
						}
						job.State = "in process"
						job.StartedAt = time.Now()
						job.CurrentValue += job.Task.D
						job.CurrentI++

						time.Sleep(time.Duration(job.Task.I) * time.Second)
					}
				}(job)
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func TTl(s *QueueService) {
	for {

		var job *model.Job
		for i := range s.Queue {
			if s.Queue[i].State == "end" {
				job = s.Queue[i]
				break
			}
		}
		if job != nil {
			go func(job *model.Job) {
				for {
					if strings.EqualFold(job.State, "end") {
						job.Task.TTL--
						time.Sleep(1 * time.Second)
						if job.Task.TTL <= 0 {
							log.Println("TTL ended. Queue len=" + strconv.Itoa(len(s.Queue)))
							return
						}
					}
				}
			}(job)
		}
		time.Sleep(5 * time.Second)
	}
}

type QueueService struct {
	Queue []*model.Job
	N     int
}

func NewService(n int) *QueueService {
	s := &QueueService{
		Queue: []*model.Job{},
		N:     n,
	}
	go process(s)
	go TTl(s)
	return s
}

func (s *QueueService) AddTask(d *model.Task) error {

	job := &model.Job{
		Task:         d,
		CurrentValue: d.N1,
	}
	// Не запускаем, но ставим статус, что нужно запустить
	job.State = "In Queue"
	job.StandedAt = time.Now()
	job.NumberOfQueue = len(s.Queue)
	s.Queue = append(s.Queue, job)
	log.Println("New task in queue. Len=" + strconv.Itoa(len(s.Queue)))

	//go func(job *model.Job) {
	//	for {
	//		if strings.EqualFold(job.State, "end") {
	//			job.Task.TTL--
	//			time.Sleep(1 * time.Second)
	//			if job.Task.TTL <= 0 {
	//				log.Println("TTL ended. Queue len=" + strconv.Itoa(len(s.Queue)))
	//				return
	//			}
	//		}
	//	}
	//}(job)

	return nil
}

func (s *QueueService) GetJobs() []*model.Job {
	for i, job := range s.Queue {
		if strings.EqualFold(job.State, "end") && job.Task.TTL <= 0 {
			s.Queue = append(s.Queue[:job.NumberOfQueue], s.Queue[job.NumberOfQueue+1:]...)
		}

		job.NumberOfQueue = i
	}
	return s.Queue
}
