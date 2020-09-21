package models

import (
	"time"
)

type UserModel struct {
	BaseModel
	Username         string            `sql:"column:username;unique_index" json:"username,omitempty"`
	EmailAddress     string            `sql:"column:email_address;size:512;unique_index" json:"email_address,omitempty"`
	TempEmailAddress string            `sql:"column:temp_email_address;size:512;" json:"-"`
	Metadata         *UserMetadata     `sql:"column:metadata;type:lob" json:"metadata,omitempty"`
	Credentials      []UserCredentials `sql:"foreignkey:user_id" json:"credentials,omitempty"`
	Inactive         bool              `sql:"column:inactive" json:"inactive,omitempty"`
}

func (um UserModel) TableName() string {
	return "t_user"
}

type UserCredentials struct {
	DeletableBaseModel
	UserID              uint       `sql:"column:user_id;not null" json:"-"`
	Type                uint8      `sql:"column:cred_type;auto_increment:false" json:"cred_type,omitempty"`
	Value               string     `sql:"column:value;size:2048" json:"value,omitempty"`
	FirstInvalidAttempt *time.Time `sql:"column:first_invalid_attempt" json:"first_invalid_attempt,omitempty"`
	InvalidAttemptCount uint       `sql:"column:invalid_attempt_count" json:"invalid_attempt_count,omitempty"`
	Bocked              bool       `sql:"column:blocked" json:"bocked,omitempty"`
}

func (uc UserCredentials) TableName() string {
	return "t_user_credentials"
}

type UserOTP struct {
	ID        uint       `sql:"column:id;primary_key" json:"id,omitempty"`
	CreatedAt time.Time  `sql:"column:created_at" json:"created_at,omitempty"`
	DeletedAt *time.Time `sql:"column:deleted_at;index" json:"deleted_at,omitempty"`
	ValueHash string     `sql:"column:hash_value;index" json:"value_hash,omitempty"`
	UserID    uint       `sql:"column:user_id" json:"user_id"`
}

func (o UserOTP) TableName() string {
	return "t_user_otp"
}

func (uc *UserCredentials) IncrementInvalidAttempt(maxAllowed uint, window time.Duration) (blocked bool) {
	if uc.Bocked || maxAllowed < 0 {
		return uc.Bocked
	}
	now := time.Now()
	if uc.FirstInvalidAttempt == nil || uc.FirstInvalidAttempt.Add(window).Before(now) {
		uc.FirstInvalidAttempt = &now
		uc.InvalidAttemptCount = 1
		if maxAllowed == 0 {
			uc.Bocked = true
		}
	} else {
		uc.InvalidAttemptCount += 1
		if maxAllowed < uc.InvalidAttemptCount {
			uc.Bocked = true
		}
	}
	return uc.Bocked
}
