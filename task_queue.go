package main

import (
	"errors"
	"fmt"
	"sync"
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
		Mutex          sync.Mutex
		BusyQueue      BusyQueue
		TaskQueue      TaskQueue
		PriorityQueue  PriorityQueue
		AddTaskChannel chan Task
		webhookurl     string
	}
)

func initTaskManager(webhookurl string) TaskManager {
	return TaskManager{
		webhookurl:     webhookurl,
		TaskQueue:      make(TaskQueue),
		BusyQueue:      make(BusyQueue),
		PriorityQueue:  make(PriorityQueue, 0),
		AddTaskChannel: make(chan Task),
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
		tm.AddTaskChannel <- task
		return nil
	case STOP_VM:
		log.Info().Msgf("adding %v to STOP_VM queue", task)
		if _, exists := tm.TaskQueue[task.notebook_id]; !exists {
			errorMsg := fmt.Sprintf("notebook with id %v is not running", task.notebook_id)
			return errors.New(errorMsg)
		}
		tm.AddTaskChannel <- task
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
			tm.AddTaskChannel <- createVm
		}
		if task.paragraph_id == "" {
			return errors.New("paragraph_id missing")
		}
		if task.code == "" {
			return errors.New("no code to execute")
		}

		tm.AddTaskChannel <- task
		return nil
	default:
		errorMsg := fmt.Sprintf("unknown task %v", task)
		log.Info().Msgf(errorMsg)
		return errors.New(errorMsg)
	}
}

func (tm *TaskManager) setupChannels() {
	for task := range tm.AddTaskChannel {
		tm.Mutex.Lock()

		tm.TaskQueue[task.notebook_id] = append(tm.TaskQueue[task.notebook_id], task)
		tm.PriorityQueue = append(tm.PriorityQueue, task.notebook_id)
		log.Info().Msgf("Task %v added to the queue\n", task)

		tm.Mutex.Unlock()

	}
}

// This is an eventloop
// TaskManager has a PriorityQueue, TaskQueue and BusyQueue
// We check PriorityQueue, get a notebook_id
// if notebook_id is not busy(ie, true in BusyQueue) -> We process it -> set BusyQueue[notebook_id] => true
// if notebook_id is busy, we check to next notebook_id in PriorityQueue
func (tm *TaskManager) setupEventLoop() error {
	for {

		if len(tm.PriorityQueue) == 0 {
			continue
		}

		task := Task{notebook_id: ""}
		// IIFE here for easy defer Mutex.Unlock
		// Otherwise, I have to always do "Unlock(); continue;" for the loop outside
		func() {
			tm.Mutex.Lock()
			defer tm.Mutex.Unlock()
			notebook_index := -1
			for i := 0; i < len(tm.PriorityQueue); i++ {
				current_notebook := tm.PriorityQueue[i]
				// current_notebook is already busy
				if tm.BusyQueue[current_notebook] {
					continue
				}
				notebook_index = i
				break
			}

			if notebook_index == -1 {
				// log.Info().Msg("all notebooks are busy")
				return
			}

			current_notebook := tm.PriorityQueue[notebook_index]
			tm.PriorityQueue = append(tm.PriorityQueue[:notebook_index], tm.PriorityQueue[notebook_index+1:]...)
			log.Info().Msgf("current_notebook: %v", current_notebook)

			current_notebook_queue, ok := tm.TaskQueue[current_notebook]
			if !ok {
				log.Info().Msgf("notebook: %v TaskQueue is not present", current_notebook)
				return
			}
			if len(current_notebook_queue) < 1 {
				log.Info().Msgf("notebook: %v TaskQueue is empty", current_notebook)
				return
			}
			task = current_notebook_queue[0]
			tm.TaskQueue[current_notebook] = tm.TaskQueue[current_notebook][1:]
			tm.BusyQueue[current_notebook] = true
		}()

		if task.notebook_id == "" {
			// log.Info().Msgf("Task is empty %v", task)
			continue
		}

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

// I will make firecracker here
func doStartVM(tm *TaskManager, task Task) {
	time.Sleep(5000 * time.Millisecond)
	tm.Mutex.Lock()
	tm.BusyQueue[task.notebook_id] = false
	webhookurl := tm.webhookurl
	tm.Mutex.Unlock()

	go sendToWebHook(webhookurl, task)
	log.Info().Msgf("%v created", task)
}

// we shouldnt delete rest of the tasks after deleting the VM
// imagine, if the next request is "CREATE_VM".
// Another approach is clear all requests until upcoming "CREATE_VM"
// If some one use this program for inspiration for some production software,
// keep that in mind
// If its "RUN_PARAGRAPH", maybe we start the VM?
func doStopVM(tm *TaskManager, task Task) {
	time.Sleep(500 * time.Millisecond)

	tm.Mutex.Lock()
	tm.BusyQueue[task.notebook_id] = false
	webhookurl := tm.webhookurl
	tm.Mutex.Unlock()

	go sendToWebHook(webhookurl, task)

	log.Info().Msgf("%v deleted", task)
}

func doRunParagraph(tm *TaskManager, task Task) {
	time.Sleep(2000 * time.Millisecond)

	tm.Mutex.Lock()
	tm.BusyQueue[task.notebook_id] = false
	webhookurl := tm.webhookurl
	tm.Mutex.Unlock()

	go sendToWebHook(webhookurl, task)

	log.Info().Msgf("%v RAN", task)
}
