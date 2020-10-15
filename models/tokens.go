package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

type (
	TokensModel struct {
		BaseModel
		RequestID      string         `gorm:"column:request_id;not null" json:"request_id,omitempty"`
		ACSignature    sql.NullString `gorm:"column:ac_signature;size:512;index:idx_token_ac" json:"ac_signature,omitempty"`
		ATSignature    sql.NullString `gorm:"column:at_signature;size:512;index:idx_token_at" json:"at_signature,omitempty"`
		RTSignature    sql.NullString `gorm:"column:rt_signature;size:512;index:idx_token_rt" json:"rt_signature,omitempty"`
		RTExpiry       sql.NullTime   `gorm:"column:rt_expiry" json:"rt_expiry,omitempty"`
		ATExpiry       sql.NullTime   `gorm:"column:at_expiry" json:"at_expiry,omitempty"`
		ACExpiry       sql.NullTime   `gorm:"column:ac_expiry" json:"ac_expiry,omitempty"`
		RequestProfile *SavedProfile  `gorm:"column:request_profile" json:"request_profile,omitempty"`
	}
	SavedProfile struct {
		Attributes map[string]string
	}
	JTIModel struct {
		ID     string    `gorm:"column:id;size:256;primary_key" json:"id"`
		Expiry time.Time `gorm:"column:expiry" json:"expiry"`
	}
)

func (tm JTIModel) AutoMigrate(db gorm.Migrator) error {
	return db.AutoMigrate(&tm)
}

func (tm TokensModel) AutoMigrate(db gorm.Migrator) error {
	return db.AutoMigrate(&tm)
}

func (sp *SavedProfile) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql", "sqlite":
		return "JSON"
	case "postgres":
		return "JSONB"
	default:
		return "lob"
	}
}

func (tm TokensModel) TableName() string {
	return "t_tokens"
}

func (tm JTIModel) TableName() string {
	return "t_assertions"
}

func (sp *SavedProfile) Scan(src interface{}) error {
	if bytes, ok := src.([]byte); ok {
		return json.Unmarshal(bytes, &sp.Attributes)
	} else {
		return errors.New("failed to unmarshal value for SavedProfile")
	}
}

func (sp *SavedProfile) Value() (marshal driver.Value, err error) {
	marshal, err = json.Marshal(sp.Attributes)
	return
}

func (sp *SavedProfile) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &sp.Attributes)
}

func (sp *SavedProfile) MarshalJSON() ([]byte, error) {
	return json.Marshal(sp.Attributes)
}
