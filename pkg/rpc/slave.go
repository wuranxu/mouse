package rpc

import (
	"context"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/wuranxu/mouse/pkg/rpc/proto"
	"github.com/wuranxu/mouse/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

const HeartbeatInterval = 5 * time.Second

// Slave slave
type Slave struct {
	client proto.MouseService_DoClient
	NodeId string
	// mac addr for machine
	mac string
	// machine ip
	ip     string
	ticker *time.Ticker
}

func (s *Slave) Recv() (*proto.Message, error) {
	return s.client.Recv()
}

func (s *Slave) Send(msgType MsgType, data []byte) error {
	return s.client.Send(NewMessage(msgType, s.NodeId, data))
}

func (s *Slave) Close() error {
	return s.client.Send(NewMessage(Quit, s.NodeId, nil))
}

func (s *Slave) Stop() {
	s.ticker.Stop()
}

// UploadStats upload cpu/memory usage and
func (s *Slave) UploadStats() error {
	for {
		select {
		case <-s.ticker.C:
			if err := s.uploadStats(); err != nil {
				log.Println("send heartbeat error: ", err)
			}
		}
	}
}

func (s *Slave) uploadStats() error {
	mem, cpu := utils.GetCpuAndMemoryUsage()
	pack, err := msgpack.Marshal(WorkNodeStat{
		State:              Ready,
		CurrentCpuUsage:    cpu,
		CurrentMemoryUsage: mem,
		Count:              0,
	})
	if err != nil {
		return err
	}
	return s.Send(Heartbeat, pack)
}

func newSlave(ctx context.Context, addr string, nodeId string) (*Slave, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := proto.NewMouseServiceClient(conn)
	cli, err := client.Do(ctx)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// Send clientReady message to server
	s := &Slave{NodeId: nodeId, client: cli, ticker: time.NewTicker(HeartbeatInterval)}
	go s.UploadStats()
	if err = s.Send(ClientReady, nil); err != nil {
		return nil, err
	}
	return s, nil
}

func NewSlave(addr string, nodeId string) (*Slave, error) {
	return newSlave(context.Background(), addr, nodeId)
}

func NewSlaveContext(context context.Context, addr, nodeId string) (*Slave, error) {
	return newSlave(context, addr, nodeId)
}
