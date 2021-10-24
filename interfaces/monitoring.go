package interfaces

import "time"

type WrappMonitoring interface {
	SendData(event, appName, message string, dateTime time.Time, dateTimeLayout string) error
}

type MonitoringAdapterInterface interface {
	Initializate() error
	SendData(message string) error
}
