package rpc

import (
	"context"
	"github.com/wuranxu/mouse/pkg/rpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type Slave struct {
	client proto.MouseService_DoClient
	NodeId string
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

func NewSlave(addr string, nodeId string) (*Slave, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := proto.NewMouseServiceClient(conn)
	cli, err := client.Do(context.Background())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	slave := &Slave{NodeId: nodeId, client: cli}
	return slave, nil
}
