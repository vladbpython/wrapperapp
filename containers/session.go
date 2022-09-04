package containers

import (
	"sync"

	"github.com/vladbpython/wrapperapp/system"
)

type SessionContainer struct {
	container map[string]*system.Session
	sync.RWMutex
}

func (c *SessionContainer) Add(key string, session *system.Session) {
	c.Lock()
	defer c.Unlock()
	c.container[key] = session
}

func (c *SessionContainer) Get(key string) *system.Session {
	c.RLock()
	defer c.RUnlock()
	return c.container[key]
}

func (c *SessionContainer) Remove(key string) {
	c.Lock()
	defer c.Unlock()
	delete(c.container, key)
}

func (c *SessionContainer) GetAll() map[string]*system.Session {
	c.RLock()
	defer c.RUnlock()
	return c.container
}

func NewSessionContainer() *SessionContainer {
	return &SessionContainer{
		container: make(map[string]*system.Session),
	}
}
