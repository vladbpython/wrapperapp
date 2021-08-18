package tools

import (
	"time"

	"github.com/vladbpython/wrapperapp/interfaces"
)

func WrapFuncOnErrorFatal(f func() error, appName string, logger interfaces.WrapLoggerInterFace, maxReties uint, timeWait uint) {
	err := f()
	if err == nil {
		return
	}
	if maxReties == 0 {
		logger.FatalError(appName, err)
		return
	}
	logger.Error(appName, err)

	for i := 1; i <= int(maxReties); i++ {
		err = f()
		if err == nil {
			return
		}
		if i == int(maxReties) {
			logger.FatalError(appName, err)
			return
		}
		logger.Error(appName, err)

		time.Sleep(time.Duration(timeWait) * time.Second)
	}

}
