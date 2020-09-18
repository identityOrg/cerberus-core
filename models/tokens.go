package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type (
	TokensModel struct {
		BaseModel
		RequestID      string        `sql:"column:request_id;not null" json:"request_id,omitempty"`
		ACSignature    string        `sql:"column:ac_signature;size:512;unique_index" json:"ac_signature,omitempty"`
		ATSignature    string        `sql:"column:at_signature;size:512;unique_index" json:"at_signature,omitempty"`
		RTSignature    string        `sql:"column:rt_signature;size:512;unique_index" json:"rt_signature,omitempty"`
		RTExpiry       time.Time     `sql:"column:rt_expiry" json:"rt_expiry,omitempty"`
		ATExpiry       time.Time     `sql:"column:at_expiry" json:"at_expiry,omitempty"`
		ACExpiry       time.Time     `sql:"column:ac_expiry" json:"ac_expiry,omitempty"`
		RequestProfile *SavedProfile `sql:"column:request_profile;type:lob" json:"request_profile,omitempty"`
	}
	SavedProfile struct {
		Attributes map[string]string
	}
	JTIModel struct {
		ID     string    `sql:"column:id;size:256;primary_key" json:"id"`
		Expiry time.Time `sql:"column:expiry" json:"expiry"`
	}
)

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
