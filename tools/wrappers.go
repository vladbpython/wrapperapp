package tools

import (
	"fmt"
	"reflect"
	"time"

	"github.com/vladbpython/wrapperapp/interfaces"
)

func WrapFunc(fn interface{}, arguments ...interface{}) ([]reflect.Value, error) {
	var values []reflect.Value
	fParser := reflect.ValueOf(fn)

	if len(arguments) != fParser.Type().NumIn() {
		return values, fmt.Errorf("func the number of arguments is not adapted.")
	}
	in := make([]reflect.Value, len(arguments))
	for k, v := range arguments {
		in[k] = reflect.ValueOf(v)
	}
	values = fParser.Call(in)
	return values, nil

}

func WrapFuncError(fn interface{}, arguments ...interface{}) error {

	values, err := WrapFunc(fn, arguments...)
	if err != nil {
		return err
	}

	for _, vObj := range values {
		value := vObj.Interface()
		switch value.(type) {
		case error:
			return value.(error)
		}
	}

	return nil
}

func WrapFuncOnErrorFatal(appName string, logger interfaces.WrapLoggerInterFace, maxReties uint, timeWait uint, fn interface{}, argumetns ...interface{}) {
	f := WrapFuncError(fn, argumetns...)
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
