package scheduler

import (
	"context"
	"log"
	"sync"
	"time"
)

type Task struct {
	ID        int
	Name      string
	TaskFunc  func()
	StartTime time.Time
	Interval  time.Duration
}

type Scheduler struct {
	tasks    []*Task
	taskChan chan *Task
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
	mu       sync.Mutex
}

func NewScheduler() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Scheduler{
		taskChan: make(chan *Task),
		ctx:      ctx,
		cancel:   cancel,
	}
	go s.worker()
	return s
}

func (s *Scheduler) AddTask(task *Task) {
	s.mu.Lock()
	s.tasks = append(s.tasks, task)
	s.mu.Unlock()
	s.taskChan <- task
}

func (s *Scheduler) worker() {
	for {
		select {
		case task := <-s.taskChan:
			s.wg.Add(1)
			go s.runTask(task)
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *Scheduler) runTask(task *Task) {
	defer s.wg.Done()
	if task.Interval == 0 {
		time.Sleep(time.Until(task.StartTime))
		s.executeTask(task)
	} else {
		ticker := time.NewTicker(task.Interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.executeTask(task)
			case <-s.ctx.Done():
				return
			}
		}
	}
}

func (s *Scheduler) executeTask(task *Task) {
	start := time.Now()
	task.TaskFunc()
	duration := time.Since(start)
	logTaskExecution(task.ID, task.Name, duration)
}

func logTaskExecution(taskID int, taskName string, duration time.Duration) {
	cyan := "\033[36m"
	purple := "\033[35m"
	green := "\033[32m"
	reset := "\033[0m"

	log.Printf("%s\"%s[TaskRunner]%s Task #%d: %s\" %sin %v%s", cyan, purple, cyan, taskID, taskName, green, duration, reset)
}

func (s *Scheduler) Stop() {
	s.cancel()
	s.wg.Wait()
}
