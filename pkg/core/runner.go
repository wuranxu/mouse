package core

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Mode string

const (
	// CommandMode started by command line
	CommandMode Mode = "command"
	// WebMode started by mouse web
	WebMode Mode = "web"
)

type IRunner interface {
	Run()
	CurrentMode() Mode
}

type RunnerOptions func(*Runner)

type Runner struct {
	ctx   context.Context
	wg    sync.WaitGroup
	env   *Environment
	mode  Mode
	tasks []*Task
}

func (r *Runner) MaxQps() int {
	return r.env.maxQps
}

func (r *Runner) CurrentMode() Mode {
	return r.mode
}

func (r *Runner) Run() {
	if r.tasks == nil || len(r.tasks) == 0 {
		return
	}
	if r.MaxQps() == 0 {
		// no qps limit
		for _, task := range r.tasks {
			task.Run()
		}
		return
	}
	for {
		select {
		case <-r.ctx.Done():
			// 停止
			log.Println("context is done1")
			return
		default:
			if err := r.env.limiter.Wait(r.ctx); err != nil {
				select {
				case <-r.ctx.Done():
					return
				default:
					log.Println("rate limiter pool get task failed, ", err)
					continue
				}
			}
			if len(r.tasks) == 1 {
				r.tasks[0].Run()
				continue
			}
			random := rand.New(rand.NewSource(time.Now().UnixNano()))
			idx := random.Intn(len(r.tasks))
			r.tasks[idx].Run()
		}
	}

}

func (r *Runner) SetTask(tasks ...*Task) {
	r.tasks = tasks
}

func WithMode(mod Mode) RunnerOptions {
	return func(runner *Runner) {
		runner.mode = mod
	}
}

func WithTasks(tasks ...*Task) RunnerOptions {
	return func(runner *Runner) {
		runner.tasks = tasks
	}
}

func WithEnvironment(env *Environment) RunnerOptions {
	return func(runner *Runner) {
		runner.env = env
	}
}

func NewRunner(ctx context.Context, opts ...RunnerOptions) *Runner {
	runner := &Runner{ctx: ctx}
	for _, opt := range opts {
		opt(runner)
	}
	return runner
}
