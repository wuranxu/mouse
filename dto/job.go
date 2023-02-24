package dto

type StartJobDto struct {
	// client address
	Addr []string `json:"addr" validate:"required,dive"`
	// scene id
	SceneId int64 `json:"sceneId" validate:"required,gt=0"`
}
