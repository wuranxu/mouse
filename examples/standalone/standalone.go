package main

import (
	"context"
	"github.com/wuranxu/mouse/pkg/core"
	"github.com/wuranxu/mouse/pkg/protocol"
	"github.com/wuranxu/mouse/pkg/protocol/http"
	"log"
	"time"
)

type PityResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	result := time.Now().Add(10 * time.Second)
	env := core.NewEnvironment(core.WithUser(10, 3, result))
	runner := core.NewRunner(ctx,
		core.WithMode(core.CommandMode),
		core.WithEnvironment(env))
	s, err := core.NewSlaveContext(ctx, "127.0.0.1:12101", "宁愿相信爱", runner, core.WithCancelFunc(cancel))
	if err != nil {
		log.Fatal("start slave error: ", err)
	}
	// func
	job := func() {
		client := http.NewHTTPClient()
		url := `https://api.pity.fun/auth/login`
		req := http.NewRequest(url, "GET", protocol.WithBody(map[string]string{
			"username": "woody", "password": "wuranxu",
		}))
		post := client.Post(req)
		if post.Error != nil {
			log.Println(post.Error)
			s.LogFailure("http", url, post.Elapsed, post.Error.Error())
			return
		}
		var resp PityResponse
		if err := post.JSON(&resp); err != nil {
			log.Println(err)
			s.LogFailure("http", url, post.Elapsed, post.Error.Error())
			return
		}
		if resp.Code != 0 {
			log.Println(resp.Msg)
			s.LogFailure("http", url, post.Elapsed, resp.Msg)
			return
		}
		log.Println("成功了你们看 ", resp.Msg)
		s.LogSuccess("http", url, post.Elapsed, 20)
	}
	runner.SetTask(core.NewScript(job))
	s.Work()

}
