package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/wuranxu/mouse/pkg/rpc"
	"io"
	"log"
	"sync"
	"time"
)

func main() {
	master, err := rpc.NewMaster("0.0.0.0", 1201)
	if err != nil {
		log.Fatal("server start error: ", err)
	}
	// 15分钟后关闭服务器
	go master.Serve()
	time.AfterFunc(5*time.Second, func() {
		master.Broadcast(rpc.Quit, nil)
	})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			generateClient("127.0.0.1:1201")
		}(&wg)
	}
	wg.Wait()
}

func generateClient(addr string) {
	random, _ := uuid.NewRandom()
	nodeId := random.String()
	slave, err := rpc.NewSlave(addr, nodeId)
	if err != nil {
		log.Fatal("start client error: ", err)
	}
	slave.Send(rpc.ClientReady, nil)
	slave.Send(rpc.Crowd, []byte(fmt.Sprintf("绝对增长，我俞强声")))
	for {
		recv, err := slave.Recv()
		if err == io.EOF {
			continue
		}
		if err != nil {
			return
		}
		switch rpc.MsgType(recv.MsgType) {
		case rpc.Quit:
			fmt.Println("该车退了")
			return
		}
	}
}
