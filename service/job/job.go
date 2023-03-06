package job

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/wuranxu/mouse/api/v1/proto"
	"github.com/wuranxu/mouse/dto"
	"github.com/wuranxu/mouse/exception"
	"github.com/wuranxu/mouse/utils/request"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
)

const MethodName = "/MouseService/Start"

const (
	ErrConnectClientCode = iota + 10010
	ErrStartClientCode
	ErrQueryScene
)

var (
	ErrConnectClient = exception.Err("can't connect to client")
	ErrStartClient   = exception.Err("can't start client")
)

func StartJob(ctx *gin.Context) (any, error) {
	body := request.GetJson[dto.StartJobDto](ctx)
	// mock data
	data, err := ioutil.ReadFile("./test_data.yml")
	if err != nil {
		return ErrQueryScene, err
	}
	task := proto.Task{
		TaskId:   3,
		Data:     data,
		MaxQps:   &body.MaxQps,
		Interval: &body.Interval,
	}
	c := context.Background()
	for _, addr := range body.Addr {
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return ErrConnectClientCode, ErrConnectClient.New(err)
		}
		var res proto.MouseResponse
		err = conn.Invoke(c, MethodName, &task, &res)
		if err != nil {
			return ErrStartClientCode, ErrStartClient.New(err)
		}
		if res.Code != 0 {
			return ErrStartClientCode, ErrStartClient.New(res.Msg)
		}
	}
	return nil, nil
}
