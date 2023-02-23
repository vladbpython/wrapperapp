package logging

import (
	"gopkg.in/natefinch/lumberjack.v2"
)

type FileLogger struct {
	out *lumberjack.Logger
	err *lumberjack.Logger
}

func (l *FileLogger) GetOutWriter() *lumberjack.Logger {
	return l.out
}

func (l *FileLogger) GetErrWriter() *lumberjack.Logger {
	return l.err
}

func (l *FileLogger) Close() {
	l.out.Close()
	l.err.Close()
}

func NewFileLogger(dirPath string, maxSize int, maxBackups int, gzip bool) *FileLogger {
	return &FileLogger{
		out: &lumberjack.Logger{
			Filename:   dirPath + "/out.log",
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			Compress:   gzip,
		},
		err: &lumberjack.Logger{
			Filename:   dirPath + "/error.log",
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			Compress:   gzip,
		},
	}
}
