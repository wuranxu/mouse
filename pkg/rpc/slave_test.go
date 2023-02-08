package rpc

import (
	"github.com/google/uuid"
	"log"
	"testing"
)

func TestNewSlave(t *testing.T) {
	random, _ := uuid.NewRandom()
	nodeId := random.String()
	slave, err := NewSlave("127.0.0.1:1201", nodeId)
	if err != nil {
		log.Fatal("start error: ", err)
	}
	defer slave.Close()
	err = slave.Send(ClientReady, []byte("我来也"))
	if err != nil {
		log.Fatal(err)
	}
	err = slave.Send(Crowd, []byte("我正在发起请求"))
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		slave.Send(Heartbeat, nil)
	}()
}
