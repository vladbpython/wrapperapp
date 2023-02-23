package wrapperapp

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/vladbpython/wrapperapp/containers"
	"github.com/vladbpython/wrapperapp/interfaces"
	"github.com/vladbpython/wrapperapp/logging"
	"github.com/vladbpython/wrapperapp/system"
	"github.com/vladbpython/wrapperapp/taskmanager"
	"github.com/vladbpython/wrapperapp/tools"
)

const moduleName = "Application"

// Структура обертки приложений
type ApplicationWrapper struct {
	appName           string
	config            ConfigWrapper
	configFilePath    string
	system            system.System //Системная структура
	logger            *logging.Logging
	ctx               context.Context    //Текущий конктекст
	finish            context.CancelFunc // Закрытие текущего контекста
	gorutinesWaiter   sync.WaitGroup     // Группы горутин
	onStopFuncs       []func()
	useInfo           bool
	sessionLogger     *logging.Logging
	containerSessions *containers.SessionContainer
}

func (a *ApplicationWrapper) Debug() bool {
	var result bool
	if a.config.System.Debug >= 1 {
		result = true
	}
	return result
}

//Иницализируем конфиг
func (a *ApplicationWrapper) LoadConfigYaml(config interface{}) {
	tools.LoadYamlConfig(a.configFilePath, config)

}

func (a *ApplicationWrapper) ReadConfigEnv(filePath string) {
	tools.ReadEnvConfig(filePath)
}

func (a *ApplicationWrapper) LoadConfigEnv(section_prefix string, config interface{}) {
	tools.ParseEnvConfig(section_prefix, config)
}

//Деллигируем Контекст приложению
func (a *ApplicationWrapper) WrapApplicationContext(app interfaces.WrapApplicationContextInterface) {
	app.WrapContext(a.ctx)

}

//Деллигируем слушателя горутин  приложению
func (a *ApplicationWrapper) WrapApplicationWaitGroup(app interfaces.WrapApplicationWaitGroupInterface) {
	app.WrapWaitGroup(&a.gorutinesWaiter)
}

func (a *ApplicationWrapper) closeSessions() {
	for key, session := range a.containerSessions.GetAll() {
		a.containerSessions.Remove(key)
		session.Stop()
	}
}

//Посылаем остановку приложения
func (a *ApplicationWrapper) close(text string) {
	a.closeSessions()
	a.finish()
	a.gorutinesWaiter.Wait()
	for _, fn := range a.onStopFuncs {
		fn()
	}
	if a.useInfo {
		a.logger.Info(a.appName, "ShutDown")
	}
	a.clear()

}

func (a *ApplicationWrapper) GetAppName() string {
	return a.appName
}

func (a *ApplicationWrapper) GetConfig() ConfigWrapper {
	return a.config
}

func (a *ApplicationWrapper) GetContext() context.Context {
	return a.ctx
}

func (a *ApplicationWrapper) GetLogger() *logging.Logging {
	return a.logger
}

func (a *ApplicationWrapper) GetSystem() *system.System {
	return &a.system
}

func (a *ApplicationWrapper) GetWG() *sync.WaitGroup {
	return &a.gorutinesWaiter
}

func (a *ApplicationWrapper) SetOnStop(fns ...func()) {
	a.onStopFuncs = make([]func(), len(fns))

	for i, fn := range fns {
		a.onStopFuncs[i] = fn
	}
}

// Инциализуруем экзмепляр системы
func (a *ApplicationWrapper) initSystem(signals ...os.Signal) {

	a.system = system.NewSystem(a.config.System.Debug, signals...)
}

//Инциализруем экзмепляр логгирования
func (a *ApplicationWrapper) initLogger() {
	a.config.System.Logger.Debug = a.Debug()
	a.logger = logging.NewLog(a.config.System.Logger)
}

// Иницализируем контекст
func (a *ApplicationWrapper) initContext() {

	a.ctx, a.finish = tools.NewContextCancel(tools.ContextBackground())
}

func (a *ApplicationWrapper) clear() {
	a.appName = ""
	a.config = ConfigWrapper{}
	a.configFilePath = ""
	a.logger = nil
	a.system = system.System{}
	a.onStopFuncs = a.onStopFuncs[:0]
	a.useInfo = false

}

// Инициализация нового диспетчера задач
func (a *ApplicationWrapper) NewTaskManager(AppName string, Logger *logging.Logging) *taskmanager.BackgroundTaskManager {
	return taskmanager.NewBackgroundTaskManager(a.ctx, AppName, Logger, &a.gorutinesWaiter)
}

// Инциализируем компоненты системы
func (a *ApplicationWrapper) Setup(signals ...os.Signal) {
	a.initSystem(signals...)
	a.initLogger()

}

// Инциализировать новый экзмепляр логгирования
func (a *ApplicationWrapper) NewLogger(DirPath string) *logging.Logging {
	newFileConfig := logging.ConfigFileLogger{
		DirPath:   fmt.Sprintf("%s/%s", a.config.System.Logger.FileConfig.DirPath, DirPath),
		MaxSize:   a.config.System.Logger.FileConfig.MaxSize,
		MaxRotate: a.config.System.Logger.FileConfig.MaxRotate,
		Gzip:      a.config.System.Logger.FileConfig.Gzip,
	}
	newConfig := logging.Config{
		Debug:      a.config.System.Logger.Debug,
		FileMode:   a.config.System.Logger.FileMode,
		StdMode:    a.config.System.Logger.StdMode,
		FileConfig: newFileConfig,
	}
	logger := logging.NewLog(newConfig)

	return logger
}

func (a *ApplicationWrapper) NewSystemInterface() interfaces.WrapSystemInterface {

	return &a.system
}

//Включаем слушателя сигналов
func (a *ApplicationWrapper) RunContextListener() {
	var eventString string

	if a.useInfo {
		a.logger.Info(a.appName, "Run")
	}

	defer a.close(eventString)

	for {
		select {
		case <-a.system.OnExitSignal():
			eventString = "closed"
			return
		case <-a.system.OnReloadSessionSignal():
			a.ReloadSessions()
		case <-a.system.OnDieSignal():
			eventString = "losed terminated"
			return
		}
	}
}

func (a *ApplicationWrapper) RunListener() {
	if a.useInfo {
		a.logger.Info(a.appName, "Run")
	}
	a.close("closed")
}

func (a *ApplicationWrapper) InitializeSessions() {
	a.sessionLogger = a.NewLogger("sessions")
}

func (a *ApplicationWrapper) NewSession(appName string) *system.Session {
	session := system.NewSession(appName, a.sessionLogger, a.GetWG())
	a.containerSessions.Add(appName, session)
	return session
}

func (a *ApplicationWrapper) ReloadSessions() {
	for _, session := range a.containerSessions.GetAll() {
		session.SendSignalReload()
	}
}

//Новый экземпляр Wrapper
func NewWrapperApplication(ApplicationName, configType, ConfigFilePath string, useInfo bool, withContext bool, signals ...os.Signal) *ApplicationWrapper {
	config := ConfigWrapper{}
	switch configType {
	case "yaml":
		tools.LoadYamlConfig(ConfigFilePath, &config)
	case "env":
		tools.ReadEnvConfig(ConfigFilePath)
		tools.ParseEnvConfig("system", &config.System)
		tools.ParseEnvConfig("logging", &config.System.Logger)
	default:
		log.Fatal("Invalid config type")
	}

	appName := ApplicationName
	if config.System.AppName != "" {
		appName += " " + config.System.AppName
	}
	app := &ApplicationWrapper{
		appName:           appName,
		config:            config,
		configFilePath:    ConfigFilePath,
		useInfo:           useInfo,
		containerSessions: containers.NewSessionContainer(),
	}
	app.Setup(signals...)
	if withContext {
		app.initContext()
	}
	return app

}
