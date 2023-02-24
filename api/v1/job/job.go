package job

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wuranxu/mouse/api/v1/proto"
	u "github.com/wuranxu/mouse/api/v1/util"
	"github.com/wuranxu/mouse/dto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
)

const MethodName = "/MouseService/Start"

const (
	ErrConnectClientCode = iota + 10010
	ErrStartClientCode
	ErrParseParametersCode
	ErrQueryScene
)

var (
	ErrConnectClient = errors.New("can't connect to client")
)

type Api struct {
	app *gin.Engine
}

func New(app *gin.Engine) *Api {
	return &Api{app: app}
}

func (a *Api) AddRoute() {
	group := a.app.Group("/job")

	// route
	group.Handle("POST", "/run", u.Wrap(startJob))
}

func startJob(ctx *gin.Context) (any, error) {
	var body dto.StartJobDto
	if err := ctx.ShouldBindJSON(&body); err != nil {
		return ErrParseParametersCode, err
	}
	// mock data
	data, err := ioutil.ReadFile("./test_data.yml")
	if err != nil {
		return ErrQueryScene, err
	}
	var qps int32 = 5
	var interval int32 = 10
	task := proto.Task{
		TaskId:   1,
		Data:     data,
		MaxQps:   &qps,
		Interval: &interval,
	}
	c := context.Background()
	for _, addr := range body.Addr {
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return ErrConnectClientCode, u.ErrWrap(ErrConnectClient, err)
		}
		var res proto.MouseResponse
		err = conn.Invoke(c, MethodName, &task, &res)
		if err != nil {
			return ErrStartClientCode, err
		}
		if res.Code != 0 {
			return ErrStartClientCode, errors.New(res.Msg)
		}
	}
	return nil, nil
}
