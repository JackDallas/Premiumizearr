package service

import (
	"sync"
	"time"

	"github.com/jackdallas/premiumizearr/internal/config"
)

type ServiceTask struct {
	TaskName      string        `json:"task_name"`
	LastCompleted time.Time     `json:"last_completed"`
	Interval      time.Duration `json:"interval"`
	IsRunning     bool          `json:"is_running"`
	function      func()
}

type TaskRunnerService struct {
	tasks      []ServiceTask
	tasksMutex *sync.RWMutex
	config     *config.Config
}

func (TaskRunnerService) New() TaskRunnerService {
	return TaskRunnerService{
		tasks:      []ServiceTask{},
		tasksMutex: &sync.RWMutex{},
	}
}

func (manager *TaskRunnerService) Init(config *config.Config) {
	manager.config = config
}

func (manager *TaskRunnerService) AddTask(taskName string, interval time.Duration, function func()) {
	manager.tasksMutex.Lock()
	defer manager.tasksMutex.Unlock()
	manager.tasks = append(manager.tasks, ServiceTask{
		TaskName:      taskName,
		LastCompleted: time.Time{},
		Interval:      interval,
		IsRunning:     false,
		function:      function,
	})
}

func (manager *TaskRunnerService) Start() {
	go func() {
		for {
			manager.tasksMutex.Lock()
			for _, task := range manager.tasks {
				if task.IsRunning {
					continue
				}
				if time.Since(task.LastCompleted) > task.Interval {
					task.IsRunning = true
					go func(task ServiceTask) {
						task.function()
						task.LastCompleted = time.Now()
						task.IsRunning = false
					}(task)
				}
			}
			manager.tasksMutex.Unlock()
			time.Sleep(time.Millisecond * 50)
		}
	}()
}
