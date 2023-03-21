package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"
)

const FORMAT = "2006-01-02 15:04:05"

// JSONTime json时间
type JSONTime struct {
	time.Time
}

func Now() JSONTime {
	return JSONTime{Time: time.Now()}
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t JSONTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

// Value insert timestamp into mysql need this function.
func (t JSONTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof jtime.Time
func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t *JSONTime) UnmarshalJSON(p []byte) error {
	var timeOrigin string
	err := json.Unmarshal(p, &timeOrigin)
	if err != nil {
		return err
	}
	tm, err := time.Parse(FORMAT, timeOrigin)
	if err != nil {
		return err
	}
	*t = JSONTime{Time: tm}
	return nil
}

type Model struct {
	ID         uint           `gorm:"primary_key" json:"id"`
	CreatedAt  JSONTime       `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt  JSONTime       `gorm:"column:updated_at" json:"updatedAt"`
	CreateUser uint           `gorm:"type:int" json:"createUser"`
	UpdateUser uint           `gorm:"type:int" json:"updateUser"`
	DeletedAt  gorm.DeletedAt `sql:"index" json:"deletedAt"`
}

// BaseModel no creator info
type BaseModel struct {
	ID        uint           `gorm:"primary_key" json:"id"`
	CreatedAt JSONTime       `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt JSONTime       `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `sql:"index" json:"deletedAt"`
}

type IModel interface {
	Update(user uint)
	Create(user uint)
}

func (m *Model) Update(user uint) {
	m.UpdateUser = user
	m.UpdatedAt = Now()
}

func (m *Model) Create(user uint) {
	m.CreateUser = user
	m.UpdateUser = user
	m.CreatedAt = Now()
	m.UpdatedAt = Now()
}

var Tables = []any{
	MouseUser{}, MouseStep{}, MouseScene{},
}
