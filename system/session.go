package system

import (
	"context"
	"fmt"
	"sync"

	"github.com/vladbpython/wrapperapp/logging"
	"github.com/vladbpython/wrapperapp/taskmanager"
	"github.com/vladbpython/wrapperapp/tools"
)

type Session struct {
	appName        string
	logger         *logging.Logging
	taskOnStart    *taskmanager.Task
	taskOnStop     *taskmanager.Task
	taskKeyOnStart string
	taskKeyOnStop  string
	reloadChannel  chan bool
	ctx            context.Context
	finish         context.CancelFunc
	wg             sync.WaitGroup
	systemWG       *sync.WaitGroup
}

func (s *Session) initialize() {
	s.ctx, s.finish = tools.NewContextCancel(tools.ContextBackground())
}

func (s *Session) AddFnOnStart(fn interface{}, arguments ...interface{}) {
	s.taskOnStart = taskmanager.NewTask(s.taskKeyOnStart, fn, arguments...)
}

func (s *Session) AddFnOnStop(fn interface{}, arguments ...interface{}) {
	s.taskOnStop = taskmanager.NewTask(s.taskKeyOnStop, fn, arguments...)
}

func (s *Session) callTask(task *taskmanager.Task) {
	s.logger.Info(s.appName, fmt.Sprintf("task %s is started", task.GetName()))
	task.SetInWorking(true)
	err := tools.WrapFuncError(task.GetFn(), task.GetArguments()...)
	if err != nil {
		s.logger.Error(s.appName, fmt.Errorf("task %s is finished with errors: %s", task.GetName(), err.Error()))
	}
	task.SetInWorking(false)
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
	task := s.taskOnStart
	if task == nil {
		return
	} else if task.GetInWorking() {
		return
	}
	s.runTreadTask(task)
}

func (s *Session) runTaskOnStop() {
	task := s.taskOnStop
	if s.taskOnStop == nil {
		return
	} else if task.GetInWorking() {
		return
	}
	s.runTreadTask(task)
	s.wg.Wait()
}

func (s *Session) reload() {
	s.runTaskOnStop()
	s.runTaskOnStart()
}

func (s *Session) runMonitor() {
	s.systemWG.Add(1)
	defer func() {
		s.runTaskOnStop()
		s.systemWG.Done()
		close(s.reloadChannel)
	}()

	s.runTaskOnStart()
	for {
		select {
		case <-s.ctx.Done():
			return
		case signal := <-s.reloadChannel:
			if signal {
				s.reload()
			}
		}
	}
}

func (s *Session) SendSignalReload() {
	s.reloadChannel <- true
}

func (s *Session) Start() {
	go s.runMonitor()
}

func (s *Session) Stop() {
	s.finish()
}

func (s *Session) Wait() {
	s.wg.Wait()
}

func NewSession(appName string, logger *logging.Logging, wg *sync.WaitGroup) *Session {
	session := &Session{
		appName:        appName,
		logger:         logger,
		taskKeyOnStart: "on start",
		taskKeyOnStop:  "on stop",
		systemWG:       wg,
		reloadChannel:  make(chan bool),
	}
	session.initialize()
	return session
}
