package interfaces

import (
	"context"
	"sync"
)

//Интерфейс деллигирования системных компонентов приложению
type WrapApplicationInferface interface {
	WrapSystem(systemInterface WrapSystemInterface)
	WrapLogger(loggerInerface WrapLoggerInterFace)
}

//Интерфейс деллигирования конекста  приложению
type WrapApplicationContextInterface interface {
	WrapContext(context context.Context)
}

//Интерфейс деллигирования слушателя горутин приложению
type WrapApplicationWaitGroupInterface interface {
	WrapWaitGroup(group *sync.WaitGroup)
}
