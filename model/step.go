package model

type MouseStep struct {
	Model
	SceneId     uint   `gorm:"type:int;index" json:"sceneId"`
	Name        string `gorm:"type:varchar(32);not null" json:"name"`
	Headers     string `gorm:"type:text" json:"headers"`
	Method      string `gorm:"type:varchar(16);not null" json:"method"`
	Body        string `gorm:"type:text" json:"body"`
	StatusCheck bool   `gorm:"type:bool;default true" json:"statusCheck"`
	Query       string `gorm:"type:text" json:"query"`
	Timeout     int    `gorm:"type:int" json:"timeout"`
}
