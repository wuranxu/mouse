package mapper

import (
	"github.com/wuranxu/mouse/dao"
	"gorm.io/gorm"
)

type Mapper[T any] struct {
}

func (m *Mapper[T]) SelectById(id any, t *T) error {
	if err := dao.Conn.Model(&t).Find(&t, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return nil

}
