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
	osChan  chan os.Signal //Системный канал событий
	errChan chan bool      //Системный канал кртических ошибок
	Debug   uint8          //Уровень дебага
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
func (s *System) Setup() {
	signal.Notify(s.osChan, syscall.SIGTERM, syscall.SIGINT)
}

//Инциализация нового экземпляра сисетмы
func NewSystem(debug uint8) System {
	osCh := make(chan os.Signal)
	errCh := make(chan bool)
	system := System{
		osChan:  osCh,
		errChan: errCh,
		Debug:   debug,
	}
	system.Setup()
	return system
}
