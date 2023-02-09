package runner

type IRunner interface {
	Run() error
}

type Runner struct {
	IRunner
	env Environment
}

func (r *Runner) Start() error {
	return r.Run()
}

func NewRunner(env Environment, task *Task) *Runner {
	return &Runner{
		IRunner: task,
		env:     env,
	}
}
