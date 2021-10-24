package interfaces

type WrapLoggerInterFace interface {
	Info(AppName string, text string)
	Debug(AppName string, text string)
	Error(AppName string, err error)
	FatalError(AppName string, err error)
	Close()
}
