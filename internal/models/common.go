package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	Id        string `gorm:"type:char(36);primaryKey"`
	CreatedAt int64
	UpdatedAt int64
}

type BaseRepo struct {
	Db *gorm.DB
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now().Unix()
	if m.Id == "" {
		m.Id = uuid.New().String()
	}
	m.CreatedAt = now
	m.UpdatedAt = now
	return
}

func (m *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now().Unix()
	return
}
