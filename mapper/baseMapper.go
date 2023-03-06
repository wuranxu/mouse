package mapper

import (
	"github.com/wuranxu/mouse/dao"
	"gorm.io/gorm"
)

type Mapper[T any] struct {
}

func (m *Mapper[T]) SelectById(id any, t *T) error {
	if err := dao.Conn.Model(t).Find(t, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return nil

}

func (m *Mapper[T]) UpdateById(t *T) error {
	return dao.Conn.Model(t).Save(t).Error
}

func (m *Mapper[T]) Insert(t *T) error {
	return dao.Conn.Model(t).Create(t).Error
}

func (m *Mapper[T]) DeleteById(id any) error {
	var t T
	return dao.Conn.Model(&t).Delete(&t, `id = ?`, id).Error
}
