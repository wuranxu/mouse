package runner

type Environment struct {
	MaxQps int
}

type Option func(*Environment)

func NewEnvironment(options ...Option) *Environment {
	env := &Environment{MaxQps: 1}
	for _, opt := range options {
		opt(env)
	}
	return env
}

func WithQps(qps int) Option {
	return func(e *Environment) {
		e.MaxQps = qps
	}
}
