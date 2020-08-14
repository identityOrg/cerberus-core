package core

import "time"

type BaseModel struct {
	ID        uint       `sql:"column:id;primary_key" json:"id,omitempty"`
	CreatedAt time.Time  `sql:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time  `sql:"column:updated_at" json:"updated_at,omitempty"`
	DeletedAt *time.Time `sql:"column:deleted_at;index" json:"deleted_at,omitempty"`
}
