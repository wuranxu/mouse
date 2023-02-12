package core

import (
	"errors"
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/wuranxu/mouse/pkg/rpc"
	"github.com/wuranxu/mouse/pkg/rpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"sync"
)

var (
	ErrorClientNotExists = errors.New("client not found")
)

var (
	Ready    = "ready"
	Crowding = "crowding"
	Running  = "running"
	Cleanup  = "cleanup"
	Stopping = "stopping"
	Stopped  = "stopped"
	Missing  = "missing"
)

type Master struct {
	proto.UnimplementedMouseServiceServer
	lock    sync.Mutex
	clients sync.Map
	host    string
	port    int
	listen  net.Listener
	server  *grpc.Server
}

type WorkNode struct {
	server proto.MouseService_DoServer
	*WorkNodeStat
}

type WorkNodeStat struct {
	State              string  `json:"state"`
	CurrentCpuUsage    float64 `json:"current_cpu_usage"`
	CurrentMemoryUsage float64 `json:"current_memory_usage"`
	Count              int     `json:"count"`
}

func (m *Master) Broadcast(msgType rpc.MsgType, data []byte) {
	m.clients.Range(func(key, value any) bool {
		node := key.(string)
		worker := value.(*WorkNode)
		err := worker.server.Send(rpc.NewMessage(msgType, node, data))
		if err != nil {
			log.Printf("broadcast to client[%s] error: %v\n", node, string(data))
			return false
		}
		return true
	})
}

func (m *Master) SendMsgToClient(msgType rpc.MsgType, nodeId string, data []byte) error {
	value, ok := m.clients.Load(nodeId)
	node := value.(*WorkNode)
	if !ok {
		return ErrorClientNotExists
	}
	return node.server.Send(rpc.NewMessage(msgType, nodeId, data))
}

func (m *Master) updateClients(ch chan *proto.Message, srv proto.MouseService_DoServer) {
	for data := range ch {
		switch rpc.MsgType(data.MsgType) {
		case rpc.Heartbeat:
			var stat WorkNodeStat
			if err := msgpack.Unmarshal(data.Data, &stat); err != nil {
				log.Println("unmarshal message failed: ", err)
			}
			log.Printf("receiving from client[%s] cpu: %f%% mem: %fMB ", data.NodeId, stat.CurrentCpuUsage, stat.CurrentMemoryUsage)
			m.clients.Store(data.NodeId, &WorkNode{
				server: srv, WorkNodeStat: &stat,
			})
		case rpc.Quit, rpc.Stop:
			m.clients.Delete(data.NodeId)
			log.Println(data.NodeId, "has disconnected.")
		case rpc.Stats:
			var result map[string]interface{}
			if err := json.Unmarshal(data.Data, &result); err != nil {
				log.Println("unmarshal message failed: ", err)
			}
			log.Println(data.NodeId, "received data: ", result["stats"])
		case rpc.ClientReady, rpc.Reconnect:
			node := &WorkNode{
				server: srv,
			}
			m.clients.Store(data.NodeId, node)
			log.Printf("client[%s] has connected.\n", data.NodeId)
		default:
			log.Printf("receiving from client[%s]: %s\n.", data.NodeId, string(data.Data))
		}
	}
}

func (m *Master) Do(srv proto.MouseService_DoServer) error {
	ch := make(chan *proto.Message)
	go m.updateClients(ch, srv)
	// close ch avoid deadlock
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
			log.Println("receiving client message error:", err)
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
	m := &Master{host: host, port: port, server: server, listen: lis}
	proto.RegisterMouseServiceServer(server, m)
	return m, nil
}
