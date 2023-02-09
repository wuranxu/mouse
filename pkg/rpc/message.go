// Package rpc
// defines message between MouseMaster and MouseSlave
package rpc

import (
	"github.com/wuranxu/mouse/pkg/rpc/proto"
)

type MsgType string

var (
	Crowd         MsgType = "crowd"
	ClientReady   MsgType = "clientReady"
	Reconnect     MsgType = "reconnect"
	Quit          MsgType = "quit"
	Stop          MsgType = "stop"
	Exception     MsgType = "exception"
	Stats         MsgType = "stats"
	ClientStopped MsgType = "clientStopped"
	Heartbeat     MsgType = "heartbeat"
)

func NewMessage(msgType MsgType, nodeId string, data []byte) *proto.Message {
	return &proto.Message{
		Data:    data,
		NodeId:  nodeId,
		MsgType: string(msgType),
	}
}
