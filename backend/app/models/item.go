package models

import (
	"gorm.io/gorm"
	"time"
)

type Item struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	Stock     int            `gorm:"not null" json:"stock"`
	Price     int64          `gorm:"not null" json:"price"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type ItemCreateOrUpdateModel struct {
	Name  string `json:"name" validate:"required"`
	Stock int    `json:"stock" validate:"required"`
	Price int64  `json:"price" validate:"required"`
}
