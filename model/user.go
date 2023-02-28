package model

type MouseUser struct {
	Model
	Name        string   `gorm:"name:name;not null" json:"name"`
	Username    string   `gorm:"name:username;unique;not null" json:"username" binding:"required"`
	Email       string   `gorm:"name:email;unique;not null;type:varchar(48)" json:"email" binding:"required"`
	Password    string   `gorm:"name:password;type:varchar(48);not null" json:"-" binding:"required"`
	LastLoginAt JSONTime `gorm:"name:last_login_at;type:timestamp" json:"lastLoginAt"`
	Role        *int     `gorm:"name:type;type:smallint;not null;default 0;" json:"role"`
}

func (*MouseUser) TableName() string {
	return "mouse_user"
}
