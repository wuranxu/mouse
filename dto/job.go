package dto

type StartJobDto struct {
	// client address
	Addr []string `json:"addr" validate:"required,dive,required"`
	// scene id
	SceneId int64 `json:"sceneId" validate:"required,gt=0"`
	// maxQps
	MaxQps int32 `json:"maxQps"`
	// interval
	Interval int32 `json:"interval" validate:"required,gte=0"`
}
