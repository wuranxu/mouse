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
	ErrDataNotExistsCode
	ErrUpdateSceneCode
	ErrDeleteSceneCode
)

var (
	ErrQuerySceneParams = exception.Err("invalid sceneId")
	ErrQueryScene       = exception.Err("query scene failed")
	ErrCreateScene      = exception.Err("create scene failed")
	ErrListScene        = exception.Err("list scene failed")
	ErrDataNotExists    = exception.Err("data not exists")
	ErrUpdateScene      = exception.Err("update scene failed")
	ErrDeleteScene      = exception.Err("delete scene failed")
)

var sceneMapper = &mapper.Scene{}

func ListScene(ctx *gin.Context) (any, error) {
	name := ctx.Query("name")
	query := mapper.NewLambdaQuery[model.MouseScene]()
	query = query.Like("name", name).Preloads("Creator", "Updater")
	list, err := sceneMapper.SelectList(query)
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
	if err := sceneMapper.SelectById(sceneId, &scene); err != nil {
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
	if err := sceneMapper.Insert(s, user.ID); err != nil {
		return ErrCreateSceneCode, ErrCreateScene.New(err)
	}
	return s, nil
}

func UpdateScene(ctx *gin.Context) (any, error) {
	user := request.GetUser(ctx)
	data := request.GetJson[dto.SceneDto](ctx)
	scene := model.NewMouseScene()
	if err := sceneMapper.SelectById(data.ID, scene); err != nil {
		return ErrQuerySceneCode, ErrQueryScene
	}
	if scene == nil {
		return ErrDataNotExistsCode, ErrDataNotExists
	}
	scene.Name = data.Name
	scene.SceneType = data.SceneType
	scene.Steps = data.Steps
	if err := sceneMapper.UpdateById(scene, user.ID); err != nil {
		return ErrUpdateSceneCode, ErrUpdateScene.New(err)
	}
	return scene, nil
}

func DeleteScene(ctx *gin.Context) (any, error) {
	user := request.GetUser(ctx)
	sceneId := ctx.Query("sceneId")
	if err := sceneMapper.DeleteById(sceneId, user.ID); err != nil {
		return ErrDeleteSceneCode, ErrDeleteScene
	}
	return nil, nil
}
