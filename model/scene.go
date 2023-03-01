package model

type SceneType int

const (
	HTTP SceneType = 1
	GRPC SceneType = 2
)

// MouseScene scene table
type MouseScene struct {
	Model
	Name      string    `gorm:"type:varchar(48);not null" json:"name"`
	SceneType SceneType `gorm:"type:int;not null;default 1" json:"sceneType"`
}
