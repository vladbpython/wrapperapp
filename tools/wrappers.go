package tools

import (
	"fmt"
	"reflect"
	"time"

	"github.com/vladbpython/wrapperapp/interfaces"
)

func WrapFunc(fn interface{}, arguments ...interface{}) error {
	fParser := reflect.ValueOf(fn)

	if len(arguments) != fParser.Type().NumIn() {
		return fmt.Errorf("func the number of arguments is not adapted.")
	}

	in := make([]reflect.Value, len(arguments))
	for k, v := range arguments {
		in[k] = reflect.ValueOf(v)
	}
	for _, vObj := range fParser.Call(in) {
		value := reflect.ValueOf(vObj).Interface()
		switch value.(type) {
		case error:
			return value.(error)
		}
	}

	return nil
}

func WrapFuncOnErrorFatal(appName string, logger interfaces.WrapLoggerInterFace, maxReties uint, timeWait uint, fn interface{}, argumetns ...interface{}) {
	f := WrapFunc(fn, argumetns...)
	err := f
	if err == nil {
		return
	}
	if maxReties == 0 {
		logger.FatalError(appName, err)
		return
	}
	logger.Error(appName, err)

	for i := 1; i <= int(maxReties); i++ {
		err = f
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
