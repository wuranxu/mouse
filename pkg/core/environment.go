package core

import "time"

type Environment struct {
	maxQps   int
	limiter  *RateLimiter
	user     int64
	ramUp    int64
	lastTime time.Time
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

func WithUser(user, period int64, last time.Time) Option {
	return func(environment *Environment) {
		environment.user = user
		environment.ramUp = period
		environment.lastTime = last
	}
}
