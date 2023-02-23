package wrapperapp

import (
	"context"
	"sync"

	"github.com/vladbpython/wrapperapp/logging"
	"github.com/vladbpython/wrapperapp/taskmanager"
)

func NewLogger(config logging.Config) *logging.Logging {
	return logging.NewLog(config)
}

func NewTasManager(ctx context.Context, AppName string, Logger *logging.Logging, wg *sync.WaitGroup) *taskmanager.BackgroundTaskManager {
	return taskmanager.NewBackgroundTaskManager(ctx, AppName, Logger, wg)
}

func NewTask(name string, fn interface{}, arguments ...interface{}) *taskmanager.Task {
	return taskmanager.NewTask(name, fn, arguments...)
}
