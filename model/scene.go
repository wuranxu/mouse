package model

type SceneType int

const (
	HTTP SceneType = 1
	GRPC SceneType = 2
)

// MouseScene scene table
type MouseScene struct {
	*Model
	Name      string    `gorm:"type:varchar(48);not null;unique" json:"name"`
	SceneType SceneType `gorm:"type:int;not null;default 1" json:"sceneType"`
	Steps     string    `gorm:"type:text;not null" json:"steps"`
	Creator   MouseUser `json:"creator" gorm:"-"`
	Updater   MouseUser `json:"updater" gorm:"-"`
}

func (m *MouseScene) TableName() string {
	return "mouse_scene"
}

func NewMouseScene() *MouseScene {
	return &MouseScene{Model: &Model{}}
}
