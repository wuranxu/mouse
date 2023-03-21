package mapper

import (
	"fmt"
	"github.com/wuranxu/mouse/dao"
	"github.com/wuranxu/mouse/model"
	"gorm.io/gorm"
)

type LambdaQuery[T model.IModel] struct {
	query *gorm.DB
}

func NewLambdaQuery[T model.IModel]() *LambdaQuery[T] {
	return &LambdaQuery[T]{query: dao.Conn.DB}
}

func (w *LambdaQuery[T]) Log() *LambdaQuery[T] {
	w.query = w.query.Debug()
	return w
}

func (w *LambdaQuery[T]) Query(query string, where ...interface{}) *LambdaQuery[T] {
	w.query = w.query.Where(query, where...)
	return w
}

func (w *LambdaQuery[T]) In(field string, where ...interface{}) *LambdaQuery[T] {
	w.query = w.query.Where(fmt.Sprintf("`%s` in (?)", field), where...)
	return w
}

func (w *LambdaQuery[T]) Like(field string, where string) *LambdaQuery[T] {
	if where == "" {
		return w
	}
	w.query = w.query.Where(fmt.Sprintf("`%s` like ?", field), fmt.Sprintf(`%%%v%%`, where))
	return w
}

func (w *LambdaQuery[T]) Eq(field string, where ...interface{}) *LambdaQuery[T] {
	w.query = w.query.Where(fmt.Sprintf("`%s` = ?", field), where...)
	return w
}

func (w *LambdaQuery[T]) Between(field string, where ...interface{}) *LambdaQuery[T] {
	w.query = w.query.Where(fmt.Sprintf("`%s` between ? and ?", field), where...)
	return w
}

func (w *LambdaQuery[T]) DateBetween(field string, where ...interface{}) *LambdaQuery[T] {
	w.query = w.query.Where(fmt.Sprintf("`%s` >= ? and  `%s` < ?", field, field), where...)
	return w
}

func (w *LambdaQuery[T]) Preload(field string, conditions ...interface{}) *LambdaQuery[T] {
	w.query = w.query.Preload(field, conditions...)
	return w
}

func (w *LambdaQuery[T]) Preloads(fields ...string) *LambdaQuery[T] {
	for _, f := range fields {
		w.query = w.query.Preload(f)
	}
	return w
}

func (w *LambdaQuery[T]) Order(field string, asc bool) *LambdaQuery[T] {
	order := "ASC"
	if !asc {
		order = "DESC"
	}
	w.query = w.query.Order(field + " " + order)
	return w
}

func (w *LambdaQuery[T]) Page(page, size int, total *int64) *LambdaQuery[T] {
	offset := (page - 1) * size
	w.query = w.query.Count(total).Offset(offset).Limit(size)
	return w
}

func (w *LambdaQuery[T]) Join(field string, where ...interface{}) *LambdaQuery[T] {
	w.query = w.query.Joins(field, where...)
	return w
}

func (w *LambdaQuery[T]) InnerJoin(field string, where ...interface{}) *LambdaQuery[T] {
	w.query = w.query.InnerJoins(field, where...)
	return w
}
