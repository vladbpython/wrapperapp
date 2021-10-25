package wrapperapp

import (
	"github.com/vladbpython/wrapperapp/interfaces"
	"github.com/vladbpython/wrapperapp/tools"
)

type WrapperStruct struct {
	AppName    string
	Logger     interfaces.WrapLoggerInterFace
	maxRetries uint
	maxWait    uint
}

func (s *WrapperStruct) SetAppName(appName string) {
	s.AppName = appName
}

func (s *WrapperStruct) GetAppName() string {
	return s.AppName
}

func (s *WrapperStruct) SetMaxRetries(num uint) {
	if num == 0 {
		num = 3
	}
	s.maxRetries = num
}

func (s *WrapperStruct) SetMaxWait(num uint) {
	if num == 0 {
		num = 3
	}
	s.maxWait = num
}

func (s *WrapperStruct) LogDebug(message string) {
	s.Logger.Debug(s.AppName, message)
}

func (s *WrapperStruct) LogInfo(message string) {
	s.Logger.Info(s.AppName, message)
}

func (s *WrapperStruct) LogError(err error) {
	if err != nil {
		s.Logger.Error(s.AppName, err)
	}
}

func (s *WrapperStruct) LogCriticalError(err error) {
	if err != nil {
		s.Logger.FatalError(s.AppName, err)
	}
}

func (s *WrapperStruct) RetryCallMethod(method interface{}, arguments ...interface{}) {

	tools.WrapFuncOnErrorFatal(s.AppName, s.Logger, s.maxRetries, s.maxWait, method, arguments...)
}

func (s *WrapperStruct) LoadSettingFromConfig(config WrapperStructConfig) {
	s.SetMaxRetries(config.MaxRetries)
	s.SetMaxWait(config.MaxWait)
}
