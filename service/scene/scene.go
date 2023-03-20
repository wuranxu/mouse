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
	ErrListSceneCode
)

var (
	ErrQuerySceneParams = exception.Err("invalid sceneId")
	ErrQueryScene       = exception.Err("query scene failed")
	ErrCreateScene      = exception.Err("create scene failed")
	ErrListScene        = exception.Err("list scene failed")
)

func ListScene(ctx *gin.Context) (any, error) {
	name := ctx.Query("name")
	query := mapper.NewLambdaQuery[model.MouseScene]()
	if name != "" {
		query.Like("name", name)
	}

	query.Log().Preload("MouseUser.ID")
	list, err := mapper.SceneMapper.SelectList(query)
	if err != nil {
		return ErrListSceneCode, ErrListScene
	}
	return list, nil
}

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
	s := model.NewMouseScene()
	s.Name = data.Name
	s.SceneType = data.SceneType
	s.Steps = data.Steps
	if err := mapper.SceneMapper.Insert(s, user.ID); err != nil {
		return ErrCreateSceneCode, ErrCreateScene.New(err)
	}
	return s, nil
}
