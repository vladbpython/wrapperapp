package taskmanager

type Task struct {
	Name      string
	fn        interface{}
	arguments []interface{}
	inWorking bool
}

func (t *Task) GetName() string {
	return t.Name
}

func (t *Task) GetFn() interface{} {
	return t.fn
}

func (t *Task) GetArguments() []interface{} {
	return t.arguments
}

func (t *Task) GetInWorking() bool {
	return t.inWorking
}

func (t *Task) SetInWorking(status bool) {
	t.inWorking = status
}

func NewTask(name string, fn interface{}, arguments ...interface{}) *Task {
	return &Task{
		Name:      name,
		fn:        fn,
		arguments: arguments,
	}
}

func NewTaskStack(name string, fn interface{}, arguments ...interface{}) Task {
	return Task{
		Name:      name,
		fn:        fn,
		arguments: arguments,
	}
}
