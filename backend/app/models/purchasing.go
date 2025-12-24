package models

import (
	"time"

	"gorm.io/gorm"
)

type Purchasing struct {
	ID         uint               `gorm:"primarykey" json:"id"`
	Date       time.Time          `gorm:"type:date;not null" json:"date"`
	SupplierID uint               `json:"supplier_id"`
	Supplier   Supplier           `gorm:"foreignKey:SupplierID" json:"supplier"`
	UserID     uint               `json:"user_id"`
	User       User               `gorm:"foreignKey:UserID" json:"user"`
	GrandTotal int64              `gorm:"not null" json:"grand_total"`
	Details    []PurchasingDetail `gorm:"foreignKey:PurchasingID" json:"details"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
	DeletedAt  gorm.DeletedAt     `gorm:"index" json:"-"`
}

type PurchasingCreateOrUpdateModel struct {
	Date       string                    `json:"date" validate:"required"`
	SupplierID uint                      `json:"supplier_id" validate:"required"`
	UserID     uint                      `json:"user_id" validate:"required"`
	GrandTotal int64                     `json:"grand_total"`
	Details    []PurchasingDetailRequest `json:"details" validate:"required,min=1,dive"`
}

type PurchasingDetailRequest struct {
	ItemID uint `json:"item_id" validate:"required"`
	Qty    int  `json:"qty" validate:"required,min=1"`
}
