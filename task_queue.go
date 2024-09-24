package main

import (
	"errors"
	"fmt"
	"time"

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
	BusyQueue     map[string]bool
	TaskManager   struct {
		BusyQueue     BusyQueue
		TaskQueue     TaskQueue
		PriorityQueue PriorityQueue
	}
)

func initTaskManager() TaskManager {
	return TaskManager{
		TaskQueue:     make(TaskQueue),
		BusyQueue:     make(BusyQueue),
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
		tm.TaskQueue[task.notebook_id] = append(tm.TaskQueue[task.notebook_id], task)
		tm.PriorityQueue = append(tm.PriorityQueue, task.notebook_id)
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
			log.Info().Msgf("/n/n/n PriorityQueue Added: %v", tm)
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

func findNextFreeNotebookIdx(tm TaskManager) int {
	for i := 0; i < len(tm.PriorityQueue); i++ {
		current_notebook := tm.PriorityQueue[i]
		// current_notebook is already busy
		if tm.BusyQueue[current_notebook] {
			continue
		}
		return i
	}
	return -1
}

// This is an eventloop
// TaskManager has a PriorityQueue, TaskQueue and BusyQueue
// We check PriorityQueue, get a notebook_id
// if notebook_id is not busy(ie, true in BusyQueue) -> We process it -> set BusyQueue[notebook_id] => true
// if notebook_id is busy, we check to next notebook_id in PriorityQueue
func (tm *TaskManager) process() error {
	for {
		// uncomment if event loop needs to run in intervals
		// time.Sleep(10 * time.Millisecond)
		if len(tm.PriorityQueue) == 0 {
			log.Info().Msgf("/n/n appState %v", tm)
			log.Info().Msg("PriorityQueue is empty")
			continue
		}

		notebook_index := findNextFreeNotebookIdx(*tm)
		if notebook_index == -1 {
			log.Info().Msg("all notebooks are busy")
			continue
		}

		current_notebook := tm.PriorityQueue[notebook_index]
		tm.PriorityQueue = append(tm.PriorityQueue[:notebook_index], tm.PriorityQueue[notebook_index+1:]...)
		log.Info().Msgf("current_notebook: %v", current_notebook)

		current_notebook_queue, ok := tm.TaskQueue[current_notebook]
		if !ok {
			log.Info().Msgf("notebook: %v TaskQueue is not present", current_notebook)
			continue
		}
		if len(current_notebook_queue) > 0 {
			log.Info().Msgf("notebook: %v TaskQueue is empty", current_notebook)
			continue
		}
		task := current_notebook_queue[0]
		tm.TaskQueue[current_notebook] = tm.TaskQueue[current_notebook][1:]
		tm.BusyQueue[current_notebook] = true

		log.Info().Msgf("Executing task: %v", task)
		switch task.Action {
		case CREATE_VM:
			log.Info().Msgf("CREATE_VM: %v", task)
			go doStartVM(tm, task)
		case STOP_VM:
			log.Info().Msgf("STOP_VM: %v", task)
			go doStopVM(tm, task)
		case RUN_PARAGRAPH:
			log.Info().Msgf("RUN_PARAGRAPH: %v", task)
			go doRunParagraph(tm, task)
		default:
			log.Warn().Msgf("unknown task: %v", task)
		}
	}
}

func doStartVM(tm *TaskManager, task Task) {
	time.Sleep(5000)
	tm.BusyQueue[task.notebook_id] = false
	log.Info().Msgf("%v created", task)
}

func doStopVM(tm *TaskManager, task Task) {
	time.Sleep(500)
	tm.BusyQueue[task.notebook_id] = false
	log.Info().Msgf("%v deleted", task)
}

func doRunParagraph(tm *TaskManager, task Task) {
	time.Sleep(2000)
	tm.BusyQueue[task.notebook_id] = false
	log.Info().Msgf("%v RAN", task)
}
