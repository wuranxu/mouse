package slave

import (
	"context"
	"github.com/google/uuid"
	"github.com/wuranxu/mouse/pkg/core"
	"github.com/wuranxu/mouse/pkg/rpc"
	"log"
	"testing"
)

func TestNewSlave(t *testing.T) {
	random, _ := uuid.NewRandom()
	nodeId := random.String()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	slave, err := NewSlaveContext(ctx, "127.0.0.1:1201", nodeId, core.NewRunner(nil, nil))
	if err != nil {
		log.Fatal("start slave error: ", err)
	}
	err = slave.Send(rpc.ClientReady, []byte("我来也"))
	if err != nil {
		log.Fatal(err)
	}
	err = slave.Send(rpc.Crowd, []byte("我正在发起请求"))
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		slave.Send(rpc.Heartbeat, nil)
	}()
}
