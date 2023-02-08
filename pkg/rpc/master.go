package rpc

import (
	"errors"
	"fmt"
	"github.com/wuranxu/mouse/pkg/rpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"io"
	"net"
	"sync"
)

var (
	ErrorClientNotExists = errors.New("client not found")
)

type Master struct {
	proto.UnimplementedMouseServiceServer
	host    string
	port    int
	lock    sync.Mutex
	clients map[string]proto.MouseService_DoServer
	listen  net.Listener
	server  *grpc.Server
	once    sync.Once
}

func (m *Master) Broadcast(msgType MsgType, data []byte) error {
	for n, srv := range m.clients {
		err := srv.Send(NewMessage(msgType, n, data))
		if err != nil {
			fmt.Println("error: ", err)
			return err
		}
	}
	return nil
}

func (m *Master) SendMsgToClient(msgType MsgType, nodeId string, data []byte) error {
	srv, ok := m.clients[nodeId]
	if !ok {
		return ErrorClientNotExists
	}
	return srv.Send(NewMessage(msgType, nodeId, data))
}

func (m *Master) updateClients(ch chan *proto.Message, srv proto.MouseService_DoServer) {
	for data := range ch {
		switch MsgType(data.MsgType) {
		case Quit, Stop:
			delete(m.clients, data.NodeId)
			fmt.Println(data.NodeId, "has disconnected...")
		case ClientReady, Reconnect:
			m.clients[data.NodeId] = srv
			fmt.Println(data.NodeId, srv, "has connected...")
		default:
			fmt.Println(string(data.Data))
		}
	}
}

func (m *Master) Do(srv proto.MouseService_DoServer) error {
	ch := make(chan *proto.Message)
	go m.updateClients(ch, srv)
	defer close(ch)
	for {
		recv, err := srv.Recv()
		if err == io.EOF {
			continue
		}
		if err != nil {
			if s, ok := status.FromError(err); ok && s != nil {
				return nil
			}
			fmt.Println("recv error:", err)
			return err
		}
		ch <- recv
	}
}

func (m *Master) Serve() error {
	return m.server.Serve(m.listen)
}

func NewMaster(host string, port int) (*Master, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}
	server := grpc.NewServer()
	m := &Master{
		host: host, port: port,
		clients: make(map[string]proto.MouseService_DoServer),
	}
	proto.RegisterMouseServiceServer(server, m)
	m.listen = lis
	m.server = server
	return m, nil
}
