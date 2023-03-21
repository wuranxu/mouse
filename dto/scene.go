package dto

import "github.com/wuranxu/mouse/model"

type SceneDto struct {
	ID        uint            `json:"id"`
	Name      string          `json:"name" validate:"required"`
	Steps     string          `json:"steps" validate:"required"`
	SceneType model.SceneType `json:"sceneType" validate:"required"`
}
