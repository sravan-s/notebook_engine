package main

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

type ActionType string

const (
	CREATE_VM     ActionType = "CREATE_VM"
	STOP_VM       ActionType = "STOP_VM"
	RUN_PARAGRAPH ActionType = "RUN_PARAGRAPH"
)

type Task struct {
	Id           int64
	Action       ActionType
	notebook_id  string
	paragraph_id string
	code         string
}

type (
	TaskQueue     map[string][]Task
	PriorityQueue []string
	TaskManager   struct {
		TaskQueue     TaskQueue
		PriorityQueue PriorityQueue
	}
)

func initTaskManager() TaskManager {
	return TaskManager{
		TaskQueue:     make(TaskQueue),
		PriorityQueue: make(PriorityQueue, 0),
	}
}

func (tm *TaskManager) addTask(task Task) error {
	if task.notebook_id == "" {
		return errors.New("Task doesnt have notebook_id")
	}

	switch task.Action {
	case CREATE_VM:
		log.Info().Msgf("adding %v to CREATE_VM queue", task)
		if _, exists := tm.TaskQueue[task.notebook_id]; exists {
			return errors.New("vm already exists for notebook, we are not doing multiple VMs yet")
		}
		tm.TaskQueue[task.notebook_id] = append(tm.TaskQueue[task.notebook_id], task)
		tm.PriorityQueue = append(tm.PriorityQueue, task.notebook_id)
		return nil
	case STOP_VM:
		log.Info().Msgf("adding %v to STOP_VM queue", task)
		if _, exists := tm.TaskQueue[task.notebook_id]; !exists {
			errorMsg := fmt.Sprintf("notebook with id %v is not running", task.notebook_id)
			return errors.New(errorMsg)
		}
		delete(tm.TaskQueue, task.notebook_id)
		return nil
	case RUN_PARAGRAPH:
		log.Info().Msgf("adding %v to RUN_PARAGRAPH queue", task)
		if _, exists := tm.TaskQueue[task.notebook_id]; !exists {
			createVm := Task{
				Id:           0,
				Action:       "CREATE_VM",
				notebook_id:  task.notebook_id,
				paragraph_id: "",
				code:         "",
			}
			tm.TaskQueue[task.notebook_id] = append(tm.TaskQueue[task.notebook_id], createVm)
			tm.PriorityQueue = append(tm.PriorityQueue, task.notebook_id)
		}
		if task.paragraph_id == "" {
			return errors.New("paragraph_id missing")
		}
		if task.code == "" {
			return errors.New("no code to execute")
		}

		tm.PriorityQueue = append(tm.PriorityQueue, task.notebook_id)
		tm.TaskQueue[task.notebook_id] = append(tm.TaskQueue[task.notebook_id], task)
		return nil
	default:
		errorMsg := fmt.Sprintf("unknown task %v", task)
		log.Info().Msgf(errorMsg)
		return errors.New(errorMsg)
	}
}
