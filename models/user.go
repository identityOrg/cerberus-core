package models

import (
	"gorm.io/gorm"
	"time"
)

type UserModel struct {
	BaseModel
	Username         string            `gorm:"column:username;index:idx_user_name,unique" json:"username,omitempty"`
	EmailAddress     string            `gorm:"column:email_address;size:512;index:idx_user_email,unique" json:"email_address,omitempty"`
	TempEmailAddress string            `gorm:"column:temp_email_address;size:512;" json:"-"`
	Metadata         *UserMetadata     `gorm:"column:metadata" json:"metadata,omitempty"`
	Credentials      []UserCredentials `gorm:"foreignKey:UserID" json:"credentials,omitempty"`
	Inactive         bool              `gorm:"column:inactive" json:"inactive,omitempty"`
}

func (um UserModel) AutoMigrate(db gorm.Migrator) error {
	return db.AutoMigrate(&um)
}

func (um UserModel) TableName() string {
	return "t_user"
}

type UserCredentials struct {
	DeletableBaseModel
	UserID              uint       `gorm:"column:user_id;not null;index:uk_user_cred_type,unique" json:"-"`
	Type                uint8      `gorm:"column:cred_type;auto_increment:false;index:uk_user_cred_type,unique" json:"cred_type,omitempty"`
	Value               string     `gorm:"column:value;size:2048" json:"value,omitempty"`
	FirstInvalidAttempt *time.Time `gorm:"column:first_invalid_attempt" json:"first_invalid_attempt,omitempty"`
	InvalidAttemptCount uint       `gorm:"column:invalid_attempt_count" json:"invalid_attempt_count,omitempty"`
	Bocked              bool       `gorm:"column:blocked" json:"bocked,omitempty"`
}

func (uc UserCredentials) AutoMigrate(db gorm.Migrator) error {
	return db.AutoMigrate(&uc)
}

func (uc UserCredentials) TableName() string {
	return "t_user_credentials"
}

type UserOTP struct {
	ID        uint       `gorm:"column:id;primary_key" json:"id,omitempty"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at,omitempty"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
	ValueHash string     `gorm:"column:hash_value;index" json:"value_hash,omitempty"`
	UserID    uint       `gorm:"column:user_id;index:uk_user_otp_id,unique" json:"user_id"`
}

func (o UserOTP) AutoMigrate(db gorm.Migrator) error {
	return db.AutoMigrate(&o)
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
