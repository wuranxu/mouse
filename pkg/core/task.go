package core

import (
	"errors"
	"fmt"
)

var (
	CommandNotFound = errors.New("there aren't scripts for test")
)

type MouseCommand func()

type Task struct {
	// test scene id
	SceneId int64
	command MouseCommand
	Weight  int
}

func (t *Task) Run() {
	if t.SceneId != 0 {
		t.runWithScene()
		return
	}
	// run by script
	if t.command == nil {
		return
	}
	t.command()
}

func (t *Task) runWithScene() {
	fmt.Println("run with scene, parse scene data and run")
	return
}

// NewScript use command for test just like locust/boomer
func NewScript(cmd MouseCommand) *Task {
	return &Task{command: cmd}
}

// NewScene use web platform for test
func NewScene(sceneId int64) *Task {
	return &Task{SceneId: sceneId}
}
