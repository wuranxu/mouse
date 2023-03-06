package scene

import (
	"github.com/gin-gonic/gin"
	"github.com/wuranxu/mouse/dto"
	"github.com/wuranxu/mouse/exception"
	"github.com/wuranxu/mouse/mapper"
	"github.com/wuranxu/mouse/model"
	"github.com/wuranxu/mouse/utils/request"
)

const (
	ErrQuerySceneParamsCode = iota + 10030
	ErrQuerySceneCode
	ErrCreateSceneCode
)

var (
	ErrQuerySceneParams = exception.Err("invalid sceneId")
	ErrQueryScene       = exception.Err("query scene failed")
	ErrCreateScene      = exception.Err("create scene failed")
)

func QueryScene(ctx *gin.Context) (any, error) {
	sceneId := ctx.Query("sceneId")
	if sceneId == "" {
		return ErrQuerySceneParamsCode, ErrQuerySceneParams
	}
	var scene model.MouseScene
	if err := mapper.SceneMapper.SelectById(sceneId, &scene); err != nil {
		return ErrQuerySceneCode, ErrQueryScene.New(err)
	}
	if scene.ID == 0 {
		return nil, nil
	}
	return scene, nil
}

func CreateScene(ctx *gin.Context) (any, error) {
	user := request.GetUser(ctx)
	data := request.GetJson[dto.SceneDto](ctx)
	s := &model.MouseScene{
		Name:       data.Name,
		SceneType:  data.SceneType,
		Steps:      data.Steps,
		CreateUser: user.ID,
		UpdateUser: user.ID,
	}
	if err := mapper.SceneMapper.Insert(s); err != nil {
		return ErrCreateSceneCode, ErrCreateScene.New(err)
	}
	return s, nil
}
