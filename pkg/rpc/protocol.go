// Package rpc
// defines message between MouseMaster and MouseSlave
package rpc

type MsgType string

var (
	Crowd         = "crowd"
	ClientReady   = "clientReady"
	Reconnect     = "reconnect"
	Quit          = "quit"
	Stop          = "stop"
	Exception     = "exception"
	Stats         = "stats"
	ClientStopped = "clientStopped"
)

type Master struct {
	client MouseService_DoServer
	host   string
	port   int
}

func (m *Master) Do(srv MouseService_DoServer) error {
	m.client = srv
	return nil
}

func (*Master) mustEmbedUnimplementedMouseServiceServer() {}

type Slave struct {
	client MouseService_DoServer
	nodeId string
}

func (s *Slave) Do(srv MouseService_DoServer) error {
	s.client = srv
	return nil
}

func (*Slave) mustEmbedUnimplementedMouseServiceServer() {}

func NewMaster() *Master {
	return &Master{}
}

func NewSlave() *Slave {
	return &Slave{}
}
