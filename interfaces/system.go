package interfaces

import "os"

type WrapSystemInterface interface {
	Die()
	ClearMemory()
	OnExitSignal() <-chan os.Signal
	OnDieSignal() <-chan bool
	StartTrace(dirPath string) interface{ Stop() }
	StopTrace(profile interface{ Stop() })
}
