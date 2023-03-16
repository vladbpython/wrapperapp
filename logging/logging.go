package logging

import (
	"fmt"
	"io"
	"log"
	"os"
)

const DateTimeLayout = "2006-01-02 15:04:05"

type Logging struct {
	DebugMode   bool
	fileHandler *FileLogger
	LogOut      *log.Logger
	logError    *log.Logger
}

func (l *Logging) error(AppName string, err error) {
	l.logError.Printf("[ERROR]: [APPNAME]:%s  [TEXT]: %s\n\r", AppName, err)
}

func (l *Logging) info(AppName, text string) {
	l.LogOut.Printf("[INFO]: [APPNAME]:%s  [TEXT]: %s\n\r", AppName, text)
}

func (l *Logging) debug(AppName, text string) {
	if l.DebugMode && l.fileHandler != nil {
		l.fileHandler.out.Write([]byte(fmt.Sprintf("[DEBUG]: [APPNAME]:%s  [TEXT]: %s\n\r", AppName, text)))
	}
}

func (l *Logging) critical(AppName string, err error) {
	l.logError.Printf("[CRITICAL]: [APPNAME]:%s  [TEXT]: %s\n\r", AppName, err)
}

//Запись лог уровень информативный
func (l *Logging) Info(AppName, text string) {
	l.info(AppName, text)

}

//Запись лог уровень отладчика
func (l *Logging) Debug(AppName, text string) {
	l.debug(AppName, text)
}

//Запись лог уровень ошибок
func (l *Logging) Error(AppName string, err error) {
	l.error(AppName, err)
}

//Запись лог уровень кртических ошибок
func (l *Logging) FatalError(AppName string, err error) {
	l.critical(AppName, err)
	os.Exit(-1)
}

//Закрываем логгер
func (l *Logging) Close() {
	if l.fileHandler == nil {
		return
	}
	l.fileHandler.Close()
}

//Новый экземпляр логгирования
func NewLog(config Config) *Logging {
	var fileLogger *FileLogger
	outWriters := make([]io.Writer, 0)
	errorWriters := make([]io.Writer, 0)

	if config.FileMode {
		fileLogger = NewFileLogger(config.FileConfig.DirPath, config.FileConfig.MaxSize, config.FileConfig.MaxRotate, config.FileConfig.Gzip)
		outWriters = append(outWriters, fileLogger.out)
		errorWriters = append(errorWriters, fileLogger.err)
	}

	if config.StdMode {
		outWriters = append(outWriters, os.Stdout)
		errorWriters = append(errorWriters, os.Stderr)
	}

	logInfoWriter := io.MultiWriter(outWriters...)
	logErrorWriter := io.MultiWriter(errorWriters...)

	return &Logging{
		DebugMode:   config.Debug,
		fileHandler: fileLogger,
		LogOut:      log.New(logInfoWriter, "", log.Ldate|log.Ltime),
		logError:    log.New(logErrorWriter, "", log.Ldate|log.Ltime),
	}
}
