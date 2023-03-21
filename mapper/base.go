package mapper

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/wuranxu/mouse/dao"
	"github.com/wuranxu/mouse/model"
	"gorm.io/gorm"
)

var (
	ExistsError = errors.New("data exists")
)

type Mapper[T model.IModel] struct {
}

func (m *Mapper[T]) SelectList(query *LambdaQuery[T]) ([]*T, error) {
	out := make([]*T, 0)
	err := query.query.Find(&out).Error
	if err := wrapperErr(err); err != nil {
		return []*T{}, err
	}
	return out, nil

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

func (m *Mapper[T]) UpdateById(t *T, user uint) error {
	(*t).Update(user)
	return dao.Conn.Model(t).Save(t).Error
}

func (m *Mapper[T]) Insert(t *T, user uint) error {
	(*t).Create(user)
	return dao.Conn.Model(t).Create(t).Error
}

func (m *Mapper[T]) DeleteById(id any, user uint) error {
	var t T
	return dao.Conn.Model(&t).Where(`id = ?`, id).Update("update_user", user).Delete(&t).Error
}

func wrapperErr(err error) error {
	if err != nil && err != gorm.ErrRecordNotFound {
		sqlError, ok := err.(*mysql.MySQLError)
		if ok && sqlError.Number == 1062 {
			return ExistsError
		}
		return err
	}
	return nil
}
