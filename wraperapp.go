package wrapperapp

import (
	"context"
	"sync"
	"time"

	"github.com/vladbpython/wrapperapp/interfaces"
	"github.com/vladbpython/wrapperapp/logging"
	"github.com/vladbpython/wrapperapp/monitoring"
	"github.com/vladbpython/wrapperapp/system"
	"github.com/vladbpython/wrapperapp/taskmanager"
	"github.com/vladbpython/wrapperapp/tools"
)

const moduleName = "Application"

// Структура обертки приложений
type ApplicationWrapper struct {
	AppName         string
	Config          ConfigWrapper
	ConfigFilePath  string
	System          system.System //Системная структура
	Logger          *logging.Logging
	Monitoring      *monitoring.Monitoring
	Ctx             context.Context    //Текущий конктекст
	finish          context.CancelFunc // Закрытие текущего контекста
	GorutinesWaiter sync.WaitGroup     // Группы горутин
	UseInfo         bool
}

//Иницализируем конфиг
func (a *ApplicationWrapper) InitConfig(config interface{}) {
	tools.LoadYamlConfig(a.ConfigFilePath, config)

}

//Добавляем системные компоненты приложению
func (a *ApplicationWrapper) WrapApplication(app interfaces.WrapApplicationInferface) {
	app.WrapLogger(a.Logger)
	app.WrapSystem(&a.System)

}

//Деллигируем Контекст приложению
func (a *ApplicationWrapper) WrapApplicationContext(app interfaces.WrapApplicationContextInterface) {
	app.WrapContext(a.Ctx)

}

//Деллигируем слушателя горутин  приложению
func (a *ApplicationWrapper) WrapApplicationWaitGroup(app interfaces.WrapApplicationWaitGroupInterface) {
	app.WrapWaitGroup(&a.GorutinesWaiter)
}

//Посылаем остановку приложения
func (a *ApplicationWrapper) Close(text string) {
	a.finish()
	a.GorutinesWaiter.Wait()
	if a.UseInfo {
		a.Logger.Info(a.AppName, "ShutDown")
	}

}

// Инциализуруем экзмепляр системы
func (a *ApplicationWrapper) InitSystem() {

	a.System = system.NewSystem(a.Config.System.Debug)
}

//Инциализруем экзмепляр логгирования
func (a *ApplicationWrapper) InitLogger() {
	a.Logger = logging.NewLog(
		a.Config.System.Debug,
		a.Config.System.Logger.DirPath,
		a.Config.System.Logger.MaxSize,
		a.Config.System.Logger.MaxRotate,
		a.Config.System.Logger.Gzip,
		a.Config.System.Logger.StdMode,
	)
}

func (a *ApplicationWrapper) InitMonitoring() {
	if a.Config.System.Monitoring.Use {
		monitoring, err := monitoring.NewMonitoringFromConfig(a.AppName, &a.Config.System.Monitoring)
		if err != nil {
			a.Logger.FatalError(a.AppName, err)
		}
		a.Monitoring = monitoring
		a.Logger.Info(a.AppName, "monitoring initializated successfully")
		a.Logger.SetMonitoring(a.Monitoring)
	}

}

// Иницализируем контекст
func (a *ApplicationWrapper) InitContext() {

	a.Ctx, a.finish = context.WithCancel(context.Background())
}

// Инициализация нового диспетчера задач
func (a *ApplicationWrapper) NewTaskManager(AppName string, Logger *logging.Logging) *taskmanager.BackgroundTaskManager {
	return taskmanager.NewBackgroundTaskManager(a.Ctx, AppName, Logger, &a.GorutinesWaiter)
}

// Инициализация новой задачи
func (a *ApplicationWrapper) NewTask(name string, fn interface{}, arguments ...interface{}) *taskmanager.Task {
	return taskmanager.NewTask(name, fn, arguments...)
}

// Инциализируем компоненты системы
func (a *ApplicationWrapper) Setup() {
	a.InitSystem()
	a.InitLogger()
	a.InitMonitoring()

}

// Инциализировать новый экзмепляр логгирования
func (a *ApplicationWrapper) NewLoggerInterface(DirPath string, MaxSize, MaxRotate int, Gzip bool) *logging.Logging {
	return logging.NewLog(
		a.Config.System.Debug,
		DirPath,
		MaxSize,
		MaxRotate,
		Gzip,
		a.Config.System.Logger.StdMode,
	)
}

func (a *ApplicationWrapper) NewSystemInterface() interfaces.WrapSystemInterface {

	return &a.System
}

//Включаем слушателя сигналов
func (a *ApplicationWrapper) RunContextListener() {
	if a.UseInfo {
		a.Logger.Info(a.AppName, "Run")
	}

	for {
		select {
		case <-a.System.OnExitSignal():
			a.Close("closed")
			return
		case <-a.System.OnDieSignal():
			a.Close("closed terminated")
			return
		case <-time.After(1 * time.Second):
			continue
		}

	}

}

func (a *ApplicationWrapper) RunListener() {
	a.GorutinesWaiter.Wait()
}

//Новый экземпляр Wrapper
func InitWrapperApplication(ApplicationName string, ConfigFilePath string, useInfo bool, withContext bool) *ApplicationWrapper {
	config := ConfigWrapper{}
	tools.LoadYamlConfig(ConfigFilePath, &config)
	appName := ApplicationName
	if config.System.AppName != "" {
		appName += " " + config.System.AppName
	}
	app := &ApplicationWrapper{
		AppName:        appName,
		Config:         config,
		ConfigFilePath: ConfigFilePath,
		UseInfo:        useInfo,
	}
	app.Setup()
	if withContext {
		app.InitContext()
	}
	return app

}
