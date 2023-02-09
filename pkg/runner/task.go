package runner

import (
	"errors"
	"fmt"
)

var (
	CommandNotFound = errors.New("there aren't scripts for test")
)

type MouseCommand func() error

type Task struct {
	// test scene id
	SceneId int64
	command MouseCommand
}

func (t *Task) Run() error {
	if t.SceneId == 0 {
		return t.runWithScene()
	}
	// run by script
	if t.command == nil {
		return CommandNotFound
	}
	return t.command()
}

func (t *Task) runWithScene() error {
	fmt.Println("run with scene, parse scene data and run")
	return nil
}

// NewScript use command for test just like locust/boomer
func NewScript(cmd MouseCommand) *Task {
	return &Task{command: cmd}
}

// NewScene use web platform for test
func NewScene(sceneId int64) *Task {
	return &Task{SceneId: sceneId}
}
