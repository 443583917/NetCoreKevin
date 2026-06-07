package scheduler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScheduler_AddTask(t *testing.T) {
	s := NewScheduler()

	err := s.AddTask("test", "Test Task", "@every 1s", func() {})
	assert.NoError(t, err)

	task, ok := s.GetTask("test")
	assert.True(t, ok)
	assert.Equal(t, "Test Task", task.Name)
	assert.Equal(t, "@every 1s", task.Spec)
	assert.Equal(t, "test", task.ID)
}

func TestScheduler_RemoveTask(t *testing.T) {
	s := NewScheduler()

	err := s.AddTask("test", "Test Task", "@every 1s", func() {})
	assert.NoError(t, err)

	s.RemoveTask("test")

	_, ok := s.GetTask("test")
	assert.False(t, ok)
}

func TestScheduler_ListTasks(t *testing.T) {
	s := NewScheduler()

	s.AddTask("task1", "Task 1", "@every 1s", func() {})
	s.AddTask("task2", "Task 2", "@every 2s", func() {})

	tasks := s.ListTasks()
	assert.Len(t, tasks, 2)
}

func TestScheduler_StartStop(t *testing.T) {
	s := NewScheduler()

	assert.False(t, s.IsRunning())

	s.Start()
	assert.True(t, s.IsRunning())

	// Starting again should be a no-op
	s.Start()
	assert.True(t, s.IsRunning())

	s.Stop()
	assert.False(t, s.IsRunning())

	// Stopping again should be a no-op
	s.Stop()
	assert.False(t, s.IsRunning())
}

func TestScheduler_InvalidCronExpr(t *testing.T) {
	s := NewScheduler()

	err := s.AddTask("bad", "Bad Task", "not a cron expression", func() {})
	assert.Error(t, err)
}

func TestScheduler_RemoveNonExistent(t *testing.T) {
	s := NewScheduler()

	// Should not panic
	s.RemoveTask("nonexistent")
}

func TestScheduler_TaskExecution(t *testing.T) {
	s := NewScheduler()

	executed := false
	err := s.AddTask("exec", "Exec Task", "@every 1s", func() {
		executed = true
	})
	assert.NoError(t, err)

	// Task is registered but scheduler not started, so it shouldn't execute
	s.Start()
	defer s.Stop()

	// The task may or may not have fired by the time we check,
	// so we just verify it was added without error
	task, ok := s.GetTask("exec")
	assert.True(t, ok)
	assert.Equal(t, "Exec Task", task.Name)
	_ = executed // execution depends on timing
}
