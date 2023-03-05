package scene

import (
	"errors"
	"github.com/gin-gonic/gin"
	u "github.com/wuranxu/mouse/api/v1/util"
	"github.com/wuranxu/mouse/mapper"
	"github.com/wuranxu/mouse/model"
)

const (
	ErrQuerySceneParamsCode = iota + 10030
	ErrQuerySceneCode
)

var (
	ErrQuerySceneParams = errors.New("invalid sceneId")
	ErrQueryScene       = errors.New("query scene error")
)

type Api struct {
	app *gin.Engine
}

func New(app *gin.Engine) *Api {
	return &Api{app: app}
}

func (a *Api) AddRoute(middleware ...gin.HandlerFunc) {
	group := a.app.Group("/scene", middleware...)
	group.GET("/", u.Wrap(queryScene))
}

func queryScene(ctx *gin.Context) (any, error) {
	sceneId := ctx.Query("sceneId")
	if sceneId == "" {
		return ErrQuerySceneParamsCode, ErrQuerySceneParams
	}
	var scene model.MouseScene
	if err := mapper.SceneMapper.SelectById(sceneId, &scene); err != nil {
		return ErrQuerySceneCode, u.ErrWrap(ErrQueryScene, err)
	}
	if scene.ID == 0 {
		return nil, nil
	}
	return scene, nil
}
