package models

import "time"

type BaseModel struct {
	DeletableBaseModel
	DeletedAt *time.Time `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
}

type DeletableBaseModel struct {
	ID        uint      `gorm:"column:id;primary_key" json:"id,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`
}
