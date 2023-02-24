package dto

type StartJobDto struct {
	// client address
	Addr []string `json:"addr" binding:"required,dive,required"`
	// scene id
	SceneId int64 `json:"sceneId" binding:"required,gt=0"`
	// maxQps
	MaxQps int32 `json:"maxQps"`
	// interval
	Interval int32 `json:"interval" binding:"required,gte=0"`
}
