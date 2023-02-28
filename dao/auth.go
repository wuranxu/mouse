package dao

import (
	"github.com/wuranxu/mouse/model"
)

func CreateUser(user *model.MouseUser) error {
	return Conn.Save(user)
}
