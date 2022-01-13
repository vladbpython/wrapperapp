package taskmanager

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/vladbpython/wrapperapp/helpers"
	loggining "github.com/vladbpython/wrapperapp/logging"
	"github.com/vladbpython/wrapperapp/tools"
)

type TaskWrapper struct {
	Task     *Task
	Type     string
	Interval time.Duration
	CallBack func(data []interface{})
	Running  bool
	Chain    chan bool
}

//Инциализируем обертку для задачи
func (t *TaskWrapper) Init() {
	t.Chain = make(chan bool)
}

//Устанавливаем/Изминяем статус обертки для задачи
func (t *TaskWrapper) ChangeStatus(status bool) {
	t.Running = status
}

//Остонвливаем задачу
func (t *TaskWrapper) StopTask() {
	t.Chain <- true
}

//Событие при остановке задачи, читаем из канала
func (t *TaskWrapper) OnStop() <-chan bool {
	return t.Chain
}

func (t *TaskWrapper) SendCallBack(data []interface{}) {
	if t.CallBack == nil {
		return
	}
	t.CallBack(data)
}

//Закрываем задачу
func (t *TaskWrapper) Clear() {
	t.ChangeStatus(false)
	close(t.Chain)
}

type BackgroundTaskManager struct {
	AppName      string
	Active       bool
	Logger       *loggining.Logging
	tasks        map[string]*TaskWrapper
	deferedTasks chan Task
	lock         sync.Mutex
	ctx          context.Context
	waitGroup    *sync.WaitGroup
}

// При ошибке пишем в лог
func (t *BackgroundTaskManager) onError(err error) {
	if err != nil {
		t.Logger.Error(t.AppName, err)
	}

}

// Событие когда задача не найдена
func (t *BackgroundTaskManager) onTaskNotFound(taskName string) {
	t.onError(fmt.Errorf("Task '%v' not found", taskName))
}

// Достаем задачу из контейнера
func (t *BackgroundTaskManager) GetTask(taskName string) *TaskWrapper {
	task, issetTask := t.tasks[taskName]
	if !issetTask {
		return nil
	}
	return task
}

//Регистрируем задачу  и добавляем её в контейнер
func (t *BackgroundTaskManager) registerTask(taskWrapper *TaskWrapper) {
	t.lock.Lock()
	defer t.lock.Unlock()
	getTask := t.GetTask(taskWrapper.Task.Name)
	if getTask == nil {
		t.tasks[taskWrapper.Task.Name] = taskWrapper
	}
}

// Подготавливаем задачу к запуску
func (t *BackgroundTaskManager) prepearTask(taskWrapper *TaskWrapper) {
	if taskWrapper.Running {
		t.onError(fmt.Errorf("Task '%v' allready running", taskWrapper.Task.Name))
		return
	}
	taskWrapper.Init()
	switch taskWrapper.Type {
	case "interval":
		t.waitGroup.Add(1)
		go t.executeTaskByInterval(taskWrapper)
	default:
		t.waitGroup.Add(1)
		go t.executeTask(taskWrapper)
	}

}

// Инициализируем задачу
func (t *BackgroundTaskManager) initTask(taskWrapper *TaskWrapper) {
	t.registerTask(taskWrapper)
	if t.Active {
		t.prepearTask(taskWrapper)
	}
}

// Добавляем задачу
func (t *BackgroundTaskManager) AddTask(task *Task) {
	taskWrapper := TaskWrapper{Task: task}
	t.initTask(&taskWrapper)
}

// Добавляем задачу по интервалу
func (t *BackgroundTaskManager) AddTaskByInterval(task *Task, interval time.Duration) {
	taskWrapper := TaskWrapper{Task: task, Type: "interval", Interval: interval}
	t.initTask(&taskWrapper)
}

//Останавливаем и удаляем задачу из контейнера
func (t *BackgroundTaskManager) RemoveTask(taskNames ...string) {
	for _, taskName := range taskNames {
		task := t.GetTask(taskName)
		if task == nil {
			t.onTaskNotFound(taskName)
		} else {
			if task.Running && task.Type == "interval" {
				task.StopTask()
			}
			delete(t.tasks, taskName)
			t.Logger.Info(t.AppName, fmt.Sprintf("Task '%v' is stopped and removed", taskName))
		}
	}

}

// Запускаем/Перезапускаем задачу вручную
func (t *BackgroundTaskManager) RunTask(taskName string) {
	task := t.GetTask(taskName)
	if task == nil {
		t.onTaskNotFound(taskName)
	} else {
		t.prepearTask(task)
	}
}

// Останавливаем задачу
func (t *BackgroundTaskManager) StopTask(taskName string) {
	task := t.GetTask(taskName)
	if task == nil {
		t.onTaskNotFound(taskName)
	} else {
		if task.Running && task.Type == "interval" {
			task.StopTask()
		} else {
			t.onError(fmt.Errorf("Task '%v' must be running and type is 'interval'", taskName))
		}
	}

}

// Запуск задачи
func (t *BackgroundTaskManager) executeTask(taskWrapper *TaskWrapper) {
	t.Logger.Info(t.AppName, fmt.Sprintf("%s task started, type: %s", taskWrapper.Task.Name, "task"))
	defer taskWrapper.Clear()
	defer t.waitGroup.Done()
	taskWrapper.ChangeStatus(true)
	data, err := t.callTask(taskWrapper.Task)
	t.onError(err)
	taskWrapper.SendCallBack(data)
	t.Logger.Info(t.AppName, fmt.Sprintf("%s task finished, type: %s", taskWrapper.Task.Name, "task"))

}

// Запуск задачи по интервалу
func (t *BackgroundTaskManager) executeTaskByInterval(taskWrapper *TaskWrapper) {
	t.Logger.Info(t.AppName, fmt.Sprintf("%s task started, type: %s", taskWrapper.Task.Name, taskWrapper.Type))
	ticker := time.NewTicker(taskWrapper.Interval)
	defer taskWrapper.Clear()
	defer ticker.Stop()
	defer t.waitGroup.Done()
	taskWrapper.ChangeStatus(true)
	for {
		select {
		case <-t.ctx.Done():
			t.Logger.Info(t.AppName, fmt.Sprintf("%s task finished, type: %s", taskWrapper.Task.Name, taskWrapper.Type))
			return
		case <-taskWrapper.OnStop():
			t.Logger.Info(t.AppName, fmt.Sprintf("%s task stopped, type: %s", taskWrapper.Task.Name, taskWrapper.Type))
			return
		case <-ticker.C:
			data, err := t.callTask(taskWrapper.Task)
			t.onError(err)
			taskWrapper.SendCallBack(data)
		}

	}

}

// Запуск отложенной задачи
func (t *BackgroundTaskManager) ExecuteDeferedTask(task Task) {
	t.Logger.Info(t.AppName, fmt.Sprintf("%s task send to channel defered task manager", task))
	t.deferedTasks <- task
}

func (t *BackgroundTaskManager) listenerDeferedTasks() {
	t.waitGroup.Add(1)
	defer t.waitGroup.Done()
	for {
		select {
		case <-t.ctx.Done():
			return
		case deferedTask := <-t.deferedTasks:
			_, err := t.callTask(&deferedTask)
			t.onError(err)
		}

	}
}

// Выполнить задачу
func (t *BackgroundTaskManager) callTask(task *Task) ([]interface{}, error) {

	values, err := tools.WrapFunc(task.fn, task.arguments)

	return helpers.SliceReflectValuesToInterfaces(values), err
}

// Запустить диспетчер задач
func (t *BackgroundTaskManager) Run() {
	if !t.Active {
		t.Active = true
		for _, task := range t.tasks {
			t.prepearTask(task)
		}
	}
}

// Остановить диспетчер задач
func (t *BackgroundTaskManager) Stop() {
	if t.Active {
		t.Active = false
		for _, task := range t.tasks {
			if task.Running && task.Type == "interval" {
				task.StopTask()
			}
		}
	}
}

// Остановить и удалить все задачи из контейнера
func (t *BackgroundTaskManager) ClearAllTasks() {
	for taskName, task := range t.tasks {
		if task.Running && task.Type == "interval" {
			task.StopTask()
		}
		delete(t.tasks, taskName)
		t.Logger.Info(t.AppName, fmt.Sprintf("Task '%v' is stopped and removed from taskmanager", taskName))
	}
}

func NewBackgroundTaskManager(ctx context.Context, appName string, logger *loggining.Logging, waitGroup *sync.WaitGroup) *BackgroundTaskManager {
	taskManager := &BackgroundTaskManager{
		AppName:      appName,
		Logger:       logger,
		tasks:        make(map[string]*TaskWrapper),
		deferedTasks: make(chan Task),
		ctx:          ctx,
		waitGroup:    waitGroup,
	}
	go taskManager.listenerDeferedTasks()
	return taskManager
}
