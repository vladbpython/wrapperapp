package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/vladbpython/wrapperapp/monitoring"
	"gopkg.in/natefinch/lumberjack.v2"
)

const DateTimeLayout = "2006-01-02 15:04:05"

type Logging struct {
	DebugMode        uint8
	LogOut           *log.Logger
	logError         *log.Logger
	logOutInstance   *lumberjack.Logger
	logErrorInstance *lumberjack.Logger
	Monitoring       *monitoring.Monitoring
}

func (l *Logging) error(AppName string, err error) {
	l.logError.Printf("[ERROR]: [APPNAME]:%s  [TEXT]: %s\n\r", AppName, err)
}

func (l *Logging) info(AppName, text string) {
	l.LogOut.Printf("[INFO]: [APPNAME]:%s  [TEXT]: %s\n\r", AppName, text)
}

func (l *Logging) debug(AppName, text string) {
	if l.DebugMode >= 1 {
		l.LogOut.Printf("[DEBUG]: [APPNAME]:%s  [TEXT]: %s\n\r", AppName, text)
	}
}

func (l *Logging) critical(AppName string, err error) {
	l.logError.Printf("[CRITICAL]: [APPNAME]:%s  [TEXT]: %s\n\r", AppName, err)
}

func (l *Logging) SetMonitoring(monitor *monitoring.Monitoring) {
	l.Monitoring = monitor
}

func (l *Logging) SendDataToMonitor(event, AppName, message string) {
	if l.Monitoring != nil {
		err := l.Monitoring.SendData(event, AppName, message, time.Now().UTC(), DateTimeLayout)
		if err != nil {
			l.error(AppName, err)
		}
	}
}

//Запись лог уровень информативный
func (l *Logging) Info(AppName, text string) {
	l.info(AppName, text)
	l.SendDataToMonitor("INFO", AppName, text)

}

//Запись лог уровень отладчика
func (l *Logging) Debug(AppName, text string) {
	l.debug(AppName, text)
	l.SendDataToMonitor("DEBUG", AppName, text)
}

//Запись лог уровень ошибок
func (l *Logging) Error(AppName string, err error) {
	l.error(AppName, err)
	l.SendDataToMonitor("ERROR", AppName, fmt.Sprintf("%s", err))
}

//Запись лог уровень кртических ошибок
func (l *Logging) FatalError(AppName string, err error) {
	l.critical(AppName, err)
	os.Exit(-1)
}

//Закрываем логгер
func (l *Logging) Close() {
	l.logOutInstance.Close()
	l.logErrorInstance.Close()
}

//Новый экземпляр логгирования
func NewLog(debug uint8, dirPath string, maxSize int, maxBackups int, gzip bool, stdMode bool) *Logging {
	var logInfoWriter io.Writer
	var logErrorWriter io.Writer

	logInfoRotator := &lumberjack.Logger{
		Filename:   dirPath + "/out.log",
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		Compress:   gzip,
	}
	logErrorRotator := &lumberjack.Logger{
		Filename:   dirPath + "/error.log",
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		Compress:   gzip,
	}

	logInfoWriter = logInfoRotator
	logErrorWriter = logErrorRotator

	if stdMode {
		logInfoWriter = io.MultiWriter(os.Stdout, logInfoWriter)
		logErrorWriter = io.MultiWriter(os.Stderr, logErrorWriter)
	}

	return &Logging{
		DebugMode:        debug,
		LogOut:           log.New(logInfoWriter, "", log.Ldate|log.Ltime),
		logOutInstance:   logInfoRotator,
		logError:         log.New(logErrorWriter, "", log.Ldate|log.Ltime),
		logErrorInstance: logErrorRotator,
	}
}
