package core

type Environment struct {
	maxQps  int
	limiter *RateLimiter
}

type Option func(*Environment)

func NewEnvironment(options ...Option) *Environment {
	env := &Environment{maxQps: 0}
	for _, opt := range options {
		opt(env)
	}
	return env
}

//func WithQps(qps int) Option {
//	return func(e *Environment) {
//		e.MaxQps = qps
//	}
//}

func WithLimiter(qps int) Option {
	return func(environment *Environment) {
		environment.maxQps = qps
		environment.limiter = NewRateLimiter(float64(qps), qps)
	}
}
