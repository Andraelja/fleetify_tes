package models

import (
    "time"
    "gorm.io/gorm"
)

type PurchasingDetail struct {
    ID           uint           `gorm:"primarykey" json:"id"`
    PurchasingID uint           `gorm:"not null" json:"purchasing_id"`
    ItemID       uint           `gorm:"not null" json:"item_id"`
    Item         Item           `gorm:"foreignKey:ItemID" json:"item"`
    Qty          int            `gorm:"not null" json:"qty"`
    Price        int64          `gorm:"not null" json:"price"` // Harga saat beli
    SubTotal     int64          `gorm:"not null" json:"sub_total"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}