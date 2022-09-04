package system

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/vladbpython/wrapperapp/logging"
	"github.com/vladbpython/wrapperapp/taskmanager"
	"github.com/vladbpython/wrapperapp/tools"
)

type Session struct {
	appName        string
	logger         *logging.Logging
	isStarted      bool
	taskKeyOnStart string
	taskKeyOnStop  string
	taskOnStart    *taskmanager.Task
	taskOnStop     *taskmanager.Task
	osChan         chan os.Signal //Системный канал событий
	ctx            context.Context
	finish         context.CancelFunc
	wg             sync.WaitGroup
	systemWG       *sync.WaitGroup
}

func (s *Session) setStarted(action bool) {
	s.isStarted = action
}

func (s *Session) initialize() {
	s.ctx, s.finish = tools.NewContextCancel(tools.ContextBackground())
}

func (s *Session) setDefaultSignal() {
	signal.Notify(s.osChan, syscall.SIGHUP)
}

func (s *Session) AddFnOnStart(fn interface{}, arguments ...interface{}) {
	s.taskOnStart = taskmanager.NewTask(s.taskKeyOnStart, fn, arguments...)
}

func (s *Session) AddFnOnStop(fn interface{}, arguments ...interface{}) {
	s.taskOnStop = taskmanager.NewTask(s.taskKeyOnStop, fn, arguments...)
}

func (s *Session) callTask(task *taskmanager.Task) {
	s.logger.Info(s.appName, fmt.Sprintf("task %s is started", task.GetName()))
	err := tools.WrapFuncError(task.GetFn(), task.GetArguments()...)
	if err != nil {
		s.logger.Error(s.appName, fmt.Errorf("task %s is finished with errors: %s", task.GetName(), err.Error()))
	}
}

func (s *Session) runTask(task *taskmanager.Task) {
	s.wg.Add(1)
	defer s.wg.Done()
	s.callTask(task)
}

func (s *Session) runTreadTask(task *taskmanager.Task) {
	go s.runTask(task)
}

func (s *Session) runTaskOnStart() {
	if s.isStarted {
		return
	} else if s.taskOnStart == nil {
		return
	}
	s.setStarted(true)
	s.runTreadTask(s.taskOnStart)
}

func (s *Session) runTaskOnStop() {
	if !s.isStarted {
		return
	} else if s.taskOnStop == nil {
		return
	}
	s.setStarted(false)
	s.runTreadTask(s.taskOnStop)
	s.wg.Wait()
}

func (s *Session) runMonitor() {
	s.systemWG.Add(1)
	defer func() {
		s.runTaskOnStop()
		s.systemWG.Done()
	}()

	s.runTaskOnStart()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.osChan:
			s.runTaskOnStop()
			s.runTaskOnStart()
		}
	}
}

func (s *Session) Start() {
	go s.runMonitor()
}

func (s *Session) Stop() {
	s.finish()
}

func NewSession(appName string, logger *logging.Logging, wg *sync.WaitGroup) *Session {
	session := &Session{
		appName:  appName,
		logger:   logger,
		systemWG: wg,
	}
	session.initialize()
	return session
}
