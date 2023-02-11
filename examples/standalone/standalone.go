package main

import (
	"context"
	"github.com/wuranxu/mouse/pkg/core"
	"github.com/wuranxu/mouse/pkg/core/slave"
	"github.com/wuranxu/mouse/pkg/protocol"
	"github.com/wuranxu/mouse/pkg/protocol/http"
	"log"
)

type PityResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func main() {
	job := func() {
		client := http.NewHTTPClient()
		req := http.NewRequest("https://api.pity.fun/auth/login", "GET", protocol.WithBody(map[string]string{
			"username": "woody", "password": "wuranxu",
		}))
		post := client.Post(req)
		if post.Error != nil {
			log.Println(post.Error)
			return
		}
		var resp PityResponse
		if err := post.JSON(&resp); err != nil {
			log.Println(err)
			return
		}
		if resp.Code != 0 {
			log.Println(resp.Msg)
			return
		}
		log.Println("success ", resp.Msg)
	}
	ctx, cancel := context.WithCancel(context.Background())
	env := core.NewEnvironment(core.WithLimiter(100))
	runner := core.NewRunner(ctx,
		core.WithMode(core.CommandMode),
		core.WithEnvironment(env),
		core.WithTasks(core.NewScript(job)))
	s, err := slave.NewSlaveContext(ctx, "127.0.0.1:12101", "宁愿相信爱", runner, slave.WithCancelFunc(cancel))
	if err != nil {
		log.Fatal("start slave error: ", err)
	}
	s.Work()

}
