package interfaces

import "github.com/vladbpython/wrapperapp/monitoring"

type WrapLoggerInterFace interface {
	Info(AppName string, text string)
	Debug(AppName string, text string)
	Error(AppName string, err error)
	FatalError(AppName string, err error)
	SetMonitoring(monitor *monitoring.Monitoring)
	Close()
}
