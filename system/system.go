package system

import (
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/pkg/profile"
)

//Структура системы
type System struct {
	osChan        chan os.Signal //Системный канал событий
	osSessionChan chan os.Signal
	errChan       chan bool //Системный канал кртических ошибок
	Debug         uint8     //Уровень дебага
}

// Послать в канал, что есть ошибка
func (s *System) Die() {
	s.errChan <- true
	return
}

func (s *System) ClearMemory() {
	debug.FreeOSMemory()
}

func (s *System) OnExitSignal() <-chan os.Signal {
	return s.osChan
}

func (s *System) OnReloadSessionSignal() <-chan os.Signal {
	return s.osSessionChan
}

func (s *System) StartTrace(dirPath string) interface{ Stop() } {
	DirPath := dirPath
	if DirPath == "" {
		DirPath = "."
	}
	if s.Debug >= 2 {
		p := profile.Start(profile.TraceProfile, profile.ProfilePath(DirPath), profile.NoShutdownHook)
		return p
	} else {
		return nil
	}

}

func (s *System) StopTrace(profile interface{ Stop() }) {
	if profile != nil {
		profile.Stop()
	}

}

func (s *System) OnDieSignal() <-chan bool {
	return s.errChan
}

// Устанавливаем сигналы
func (s *System) Setup(signals ...os.Signal) {
	sigs := make([]os.Signal, 2)
	sigs[0] = syscall.SIGTERM
	sigs[1] = syscall.SIGINT
	sigs = append(sigs, signals...)
	signal.Notify(s.osChan, sigs...)
	signal.Notify(s.osSessionChan, syscall.SIGHUP)
}

//Инциализация нового экземпляра сисетмы
func NewSystem(debug uint8, signals ...os.Signal) System {
	system := System{
		osChan:        make(chan os.Signal),
		osSessionChan: make(chan os.Signal),
		errChan:       make(chan bool),
		Debug:         debug,
	}
	system.Setup(signals...)
	return system
}
