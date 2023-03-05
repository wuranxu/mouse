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
	Steps     string    `gorm:"type:text;not null" json:"steps"`
}

func (m *MouseScene) TableName() string {
	return "mouse_scene"
}
