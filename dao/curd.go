package dao

import (
	"errors"
	"fmt"
	"github.com/wuranxu/mouse/conf"
	"github.com/wuranxu/mouse/model"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"reflect"
)

var (
	UnSupportedDatabase = errors.New("database is not supported now")
	UpdateParamsError   = errors.New("you must provide update column and value")
	StructError         = errors.New("you must provide the same struct")
)

var (
	Conn *Cursor
)

type Columns map[string]interface{}

func NewConnect(cfg conf.SqlConfig) (cur *Cursor, err error) {
	var (
		args      string
		dialector gorm.Dialector
		db        *gorm.DB
	)
	switch cfg.Name {
	case "postgres":
		// postgres sql
		args = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password)
		dialector = postgres.Open(args)
	case "mysql":
		args = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=Local",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
		dialector = mysql.Open(args)
	default:
		err = UnSupportedDatabase
		return
	}
	db, err = gorm.Open(dialector)
	if err != nil {
		return
	}
	return &Cursor{db}, nil
}

type Cursor struct {
	*gorm.DB
}

func (c *Cursor) Like(field, name string) *Cursor {
	c.DB = c.Where(fmt.Sprintf("%s like '%%%s%%'", field, name))
	return c

}

func (c *Cursor) Llike(field, name string) *Cursor {
	c.DB = c.Where(fmt.Sprintf("%s llike '%%%s'", field, name))
	return c
}
func (c *Cursor) Rlike(field, name string) *Cursor {
	c.DB = c.Where(fmt.Sprintf("%s rlike '%s%%'", field, name))
	return c
}

func (c *Cursor) Find(out interface{}, where ...interface{}) *Cursor {
	c.DB = c.DB.Find(out, where...)
	return c
}

func (c *Cursor) Order(value string) error {
	return c.DB.Order(value).Error
}

// FindPagination find with pagination
func (c *Cursor) FindPagination(page, pageSize int, out interface{}, where ...interface{}) *Cursor {
	return c.Page(page, pageSize).Find(out, where...)
}

// FindPaginationAndOrder find with order
func (c *Cursor) FindPaginationAndOrder(page, pageSize int, order string, out interface{}, where ...interface{}) (int64, error) {
	var total int64
	c.DB = c.DB.Find(out, where...).Count(&total)
	err := c.Page(page, pageSize).Order(order)
	return total, err
}

func (c *Cursor) Table(name string) *Cursor {
	c.DB = c.DB.Table(name)
	return c
}

func (c *Cursor) Select(query interface{}, args ...interface{}) *Cursor {
	c.DB = c.DB.Select(query, args...)
	return c
}

func (c *Cursor) Sql(v interface{}, sql string, params ...interface{}) error {
	return c.Raw(sql, params...).Scan(v).Error
}

func (c *Cursor) Insert(v interface{}) error {
	return c.Create(v).Error
}

func (c *Cursor) Delete(v interface{}, where ...interface{}) error {
	return c.DB.Delete(v, where...).Error
}

func (c *Cursor) Save(v interface{}) error {
	return c.DB.Save(v).Error
}

func (c *Cursor) Updates(v interface{}, attrs ...interface{}) (int64, error) {
	switch len(attrs) {
	case 0:
		return 0, UpdateParamsError
	case 1:
		switch to := attrs[0].(type) {
		// receive map and struct
		case Columns:
			c.DB = c.Model(v).Updates(to)
		default:
			var dist string
			vType := reflect.ValueOf(v).Elem().Type().String()
			if value := reflect.ValueOf(to); value.Kind() == reflect.Ptr {
				// 指针模式
				dist = value.Elem().Type().String()
			} else {
				dist = value.Type().String()
			}
			if dist != vType {
				return 0, StructError
			}
			c.DB = c.Model(v).Updates(to)
		}
	default:
		for i := 0; i < len(attrs); i += 2 {
			if i+1 >= len(attrs) {
				break
			}
			c.DB = c.Model(v).Update(attrs[i].(string), attrs[i+1])
		}
	}
	return c.DB.RowsAffected, c.DB.Error
}

func (c *Cursor) Page(current, pageSize int) *Cursor {
	c.DB = c.Offset((current - 1) * pageSize).Limit(pageSize)
	return c
}

func InitDatabase() (err error) {
	Conn, err = NewConnect(conf.Conf.Database)
	if err != nil {
		return err
	}
	for _, table := range model.Tables {
		if err = Conn.AutoMigrate(table); err != nil {
			return err
		}
	}
	return
}
