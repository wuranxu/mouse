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
	Env() *Environment
	Mode() Mode
}

type RunnerOptions func(*Runner)

type Runner struct {
	ctx   context.Context
	wg    sync.WaitGroup
	env   *Environment
	mode  Mode
	tasks []*Task
}

func (r *Runner) Env() *Environment {
	return r.env
}

func (r *Runner) Mode() Mode {
	return r.mode
}

func (r *Runner) runWithUser() {
	if r.env.user == 0 {
		log.Println("please set stress model")
		return
	}
	if r.env.ramUp == 0 {
		r.env.ramUp = 1
	}
	step := r.env.user / r.env.ramUp
	var current int64
	r.wg.Add(int(r.env.user))
	ctx, cancel := context.WithDeadline(context.Background(), r.env.lastTime)
	defer cancel()
	for current < r.env.user {
		current += step
		users := step
		if current >= r.env.user {
			users = step + r.env.user - current
		}
		for i := 1; i <= int(users); i++ {
			random := rand.New(rand.NewSource(time.Now().UnixNano()))
			idx := random.Intn(len(r.tasks))
			go func(c context.Context) {
				defer r.wg.Done()
				for {
					select {
					case <-c.Done():
						return
					default:
						r.tasks[idx].Run()
					}
				}
			}(ctx)
		}
		time.Sleep(time.Second)
	}
	r.wg.Wait()
}

func (r *Runner) runWithLimiter() {
	for {
		select {
		case <-r.ctx.Done():
			// 停止
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

func (r *Runner) Run() {
	if r.tasks == nil || len(r.tasks) == 0 {
		return
	}
	if r.env.limiter != nil {
		r.runWithLimiter()
		return
	}
	r.runWithUser()
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
