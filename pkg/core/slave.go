package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
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

type SlaveOption func(*Slave)

// Slave slave
type Slave struct {
	client proto.MouseService_DoClient
	runner IRunner
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
	err   chan error
	stats *RequestStats
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

	go s.stats.start()

	// get statistics
	go func() {
		for {
			select {
			case data := <-s.stats.messageToRunnerChan:
				marshal, err := json.Marshal(data)
				if err != nil {
					s.err <- err
					return
				}
				if err = s.Send(rpc.Stats, marshal); err != nil {
					s.err <- err
				}
			}
		}
	}()
	defer s.Quit(true)
	// start test
	s.runner.Run()
	startTime := time.Now().Unix() - s.stats.total.StartTime
	fmt.Println(s.stats.total.NumRequests, startTime, s.stats.total.NumRequests/startTime)

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

func (s *Slave) LogSuccess(requestType, name string, responseTime int64, responseLength int64) {
	s.stats.requestSuccessChan <- &requestSuccess{
		requestType:    requestType,
		name:           name,
		responseTime:   responseTime,
		responseLength: responseLength,
	}
}

func (s *Slave) LogFailure(requestType, name string, responseTime int64, exception string) {
	s.stats.requestFailureChan <- &requestFailure{
		requestType:  requestType,
		name:         name,
		responseTime: responseTime,
		error:        exception,
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
				return
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
	pack, err := msgpack.Marshal(WorkNodeStat{
		State:              Ready,
		CurrentCpuUsage:    cpu,
		CurrentMemoryUsage: mem,
		Count:              0,
	})
	if err != nil {
		return err
	}
	return s.Send(rpc.Heartbeat, pack)
}

func newSlave(ctx context.Context, addr string, nodeId string, runner IRunner, opts ...SlaveOption) (s *Slave, err error) {
	s = &Slave{
		NodeId: nodeId,
		runner: runner, ctx: ctx,
		ticker: time.NewTicker(HeartbeatInterval),
		stats:  NewRequestStats(),
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
	if runner.Mode() == CommandMode {
		signal.Notify(s.quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	}
	return
}

func NewSlaveContext(ctx context.Context, addr, nodeId string, runner IRunner, opts ...SlaveOption) (*Slave, error) {
	return newSlave(ctx, addr, nodeId, runner, opts...)
}

func WithCancelFunc(cancelFunc context.CancelFunc) SlaveOption {
	return func(slave *Slave) {
		slave.cancel = cancelFunc
	}
}
