package slave

import (
	"context"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/wuranxu/mouse/pkg/core"
	"github.com/wuranxu/mouse/pkg/core/master"
	"github.com/wuranxu/mouse/pkg/rpc"
	"github.com/wuranxu/mouse/pkg/rpc/proto"
	"github.com/wuranxu/mouse/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const HeartbeatInterval = 5 * time.Second

type Option func(*Slave)

// Slave slave
type Slave struct {
	client proto.MouseService_DoClient
	runner core.IRunner
	ctx    context.Context
	NodeId string
	// mac addr for machine
	mac string
	// machine ip
	ip     string
	ticker *time.Ticker
	// slave global cancel func
	cancel context.CancelFunc
	// quit chan
	quit chan os.Signal
	// error chan
	err chan error
}

// Quit disconnected with master
// if self is true, need send msg to master
// else close by master
func (s *Slave) Quit(self bool) {
	// close upload
	log.Println("bye~")
	s.ticker.Stop()
	if self {
		if err := s.Send(rpc.Quit, nil); err != nil {
			log.Println("send quit info to master failed, ", err)
		}
	}
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *Slave) Recv() (*proto.Message, error) {
	return s.client.Recv()
}

func (s *Slave) Send(msgType rpc.MsgType, data []byte) error {
	return s.client.Send(rpc.NewMessage(msgType, s.NodeId, data))
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

func (s *Slave) Work() {
	go func() {
		// Upload statistics info
		s.err <- s.UploadStats()
	}()
	go s.CheckSlaveStatus()
	// Send clientReady message to server
	go func() {
		s.err <- s.Send(rpc.ClientReady, nil)
	}()
	go s.Listen()

	// start test
	s.runner.Run()

}

func (s *Slave) Listen() {
	for {
		recv, err := s.Recv()
		if err == io.EOF {
			continue
		}
		if err != nil {
			log.Println("receiving from master error: ", err)
			s.err <- err
		}
		switch rpc.MsgType(recv.MsgType) {
		case rpc.Quit:
			log.Println("receiving from master: ", string(recv.Data))
			close(s.quit)
			return
		}
	}
}

// CheckSlaveStatus a goroutine for slave, check error/quit
func (s *Slave) CheckSlaveStatus() {
	for {
		select {
		case err := <-s.err:
			if err != nil {
				s.Quit(false)
				log.Fatal("exception occurred: ", s.err)
			}
		case <-s.quit:
			log.Println("quit by user")
			s.Quit(true)
			return
		case <-s.ctx.Done():
			s.Quit(false)
			return
		}
	}
}

func (s *Slave) uploadStats() error {
	mem, cpu := utils.GetCpuAndMemoryUsage()
	pack, err := msgpack.Marshal(master.WorkNodeStat{
		State:              master.Ready,
		CurrentCpuUsage:    cpu,
		CurrentMemoryUsage: mem,
		Count:              0,
	})
	if err != nil {
		return err
	}
	return s.Send(rpc.Heartbeat, pack)
}

func newSlave(ctx context.Context, addr string, nodeId string, runner core.IRunner, opts ...Option) (s *Slave, err error) {
	s = &Slave{
		NodeId: nodeId,
		runner: runner, ctx: ctx,
		ticker: time.NewTicker(HeartbeatInterval),
	}

	// load options
	for _, opt := range opts {
		opt(s)
	}

	s.ip, err = utils.GetExternalIP()
	if err != nil {
		return
	}
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	client := proto.NewMouseServiceClient(conn)
	s.client, err = client.Do(ctx)
	if err != nil {
		return
	}
	s.quit = make(chan os.Signal, 1)
	if runner.CurrentMode() == core.CommandMode {
		signal.Notify(s.quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	}
	return
}

func NewSlaveContext(ctx context.Context, addr, nodeId string, runner core.IRunner, opts ...Option) (*Slave, error) {
	return newSlave(ctx, addr, nodeId, runner, opts...)
}

func WithCancelFunc(cancelFunc context.CancelFunc) Option {
	return func(slave *Slave) {
		slave.cancel = cancelFunc
	}
}
