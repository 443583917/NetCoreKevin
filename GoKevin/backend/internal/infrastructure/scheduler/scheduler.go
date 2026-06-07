package scheduler

import (
	"log"
	"sync"

	"github.com/robfig/cron/v3"
)

// Task represents a scheduled task
type Task struct {
	ID      string
	Name    string
	Spec    string // Cron expression
	Handler func()
	EntryID cron.EntryID
}

// Scheduler manages scheduled tasks
type Scheduler struct {
	cron    *cron.Cron
	tasks   map[string]*Task
	mu      sync.RWMutex
	running bool
}

// NewScheduler creates a new scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{
		cron:  cron.New(),
		tasks: make(map[string]*Task),
	}
}

// AddTask adds a scheduled task
func (s *Scheduler) AddTask(id, name, spec string, handler func()) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entryID, err := s.cron.AddFunc(spec, handler)
	if err != nil {
		return err
	}

	s.tasks[id] = &Task{
		ID:      id,
		Name:    name,
		Spec:    spec,
		Handler: handler,
		EntryID: entryID,
	}

	log.Printf("Task added: %s (%s) - %s", id, name, spec)
	return nil
}

// RemoveTask removes a scheduled task
func (s *Scheduler) RemoveTask(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task, ok := s.tasks[id]; ok {
		s.cron.Remove(task.EntryID)
		delete(s.tasks, id)
		log.Printf("Task removed: %s", id)
	}
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		s.cron.Start()
		s.running = true
		log.Println("Scheduler started")
	}
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		s.cron.Stop()
		s.running = false
		log.Println("Scheduler stopped")
	}
}

// GetTask returns a task by ID
func (s *Scheduler) GetTask(id string) (*Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	return task, ok
}

// ListTasks returns all tasks
func (s *Scheduler) ListTasks() []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// IsRunning returns whether the scheduler is running
func (s *Scheduler) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}
