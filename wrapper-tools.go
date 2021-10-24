package wrapperapp

import (
	"context"
	"sync"

	"github.com/vladbpython/wrapperapp/logging"
	"github.com/vladbpython/wrapperapp/monitoring"
	"github.com/vladbpython/wrapperapp/monitoring/adapters"
	"github.com/vladbpython/wrapperapp/taskmanager"
)

func NewLogger(debug uint8, dirPath string, maxSize int, maxBackups int, gzip bool, stdMode bool) *logging.Logging {
	return logging.NewLog(debug, dirPath, maxSize, maxBackups, gzip, stdMode)
}

func NewTasManager(ctx context.Context, AppName string, Logger *logging.Logging, wg *sync.WaitGroup) *taskmanager.BackgroundTaskManager {
	return taskmanager.NewBackgroundTaskManager(ctx, AppName, Logger, wg)
}

func NewTask(name string, fn interface{}, arguments ...interface{}) *taskmanager.Task {
	return taskmanager.NewTask(name, fn, arguments...)
}

func NewMonitoringConfig() monitoring.ConfigMinotiring {
	return monitoring.ConfigMinotiring{}
}

func NewMonitorinAdapterConfig() adapters.ConfigAdapter {
	return adapters.ConfigAdapter{}
}

func NewMonitoring(appName string, cfg *monitoring.ConfigMinotiring) (*monitoring.Monitoring, error) {
	return monitoring.NewMonitoringFromConfig(appName, cfg)
}
