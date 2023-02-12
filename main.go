package main

import (
	"github.com/wuranxu/mouse/pkg/core"
	"log"
)

func main() {
	master, err := core.NewMaster("0.0.0.0", 12101)
	if err != nil {
		log.Fatal("server start error: ", err)
	}
	master.Serve()
	// 15分钟后关闭服务器
	//go func() {
	//	err := master.Serve()
	//	if err != nil {
	//		log.Fatal("starting Mouse master error: ", err)
	//	}
	//}()
	//time.AfterFunc(5*time.Minute, func() {
	//	master.Broadcast(rpc.Quit, nil)
	//})
	//var wg sync.WaitGroup
	//names := []string{"鸡哥", "水鬼哥", "辉总", "方总", "美式哥", "能总", "强哥", "孟总", "右总", "老虎弟弟", "仔仔", "龙哥"}
	//for _, name := range names {
	//	wg.Add(1)
	//	go func(wg *sync.WaitGroup, name string) {
	//		defer wg.Done()
	//		generateClient("127.0.0.1:12101", name)
	//	}(&wg, name)
	//}
	//wg.Wait()
}

//func generateClient(addr string, nodeId string) {
//	slave, err := core.NewSlave(addr, nodeId)
//	defer slave.Stop()
//	if err != nil {
//		log.Fatal("start client error: ", err)
//	}
//	for {
//		recv, err := slave.Recv()
//		if err == io.EOF {
//			continue
//		}
//		if err != nil {
//			log.Println("receiving from master error: ", string(recv.Data))
//			return
//		}
//		switch rpc.MsgType(recv.MsgType) {
//		case rpc.Quit:
//			log.Println("receiving from master: ", string(recv.Data))
//			return
//		}
//	}
//}
