package taskmanager

type Task struct {
	Name      string
	fn        interface{}
	arguments []interface{}
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
