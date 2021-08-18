package monitoring

type MonitoringAdapterInterface interface {
	Initializate() error
	SendData(message string) error
}
