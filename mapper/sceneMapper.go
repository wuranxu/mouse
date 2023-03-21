package mapper

import "github.com/wuranxu/mouse/model"

type Scene struct {
	Mapper[model.MouseScene]
}

var SceneMapper = &Scene{}
