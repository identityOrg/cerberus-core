package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type UserModel struct {
	BaseModel
	Username         string            `sql:"column:username;unique_index" json:"username,omitempty"`
	EmailAddress     string            `sql:"column:email_address;size:512;unique_index" json:"email_address,omitempty"`
	TempEmailAddress string            `sql:"column:temp_email_address;size:512;unique_index" json:"-"`
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
	if uc.Bocked {
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

type UserAddress struct {
	Formatted     string `json:"formatted,omitempty"`
	StreetAddress string `json:"street_address,omitempty"`
	Locality      string `json:"locality,omitempty"`
	Region        string `json:"region,omitempty"`
	PostalCode    string `json:"postal_code,omitempty"`
	Country       string `json:"country,omitempty"`
}

type UserMetadata struct {
	Name                string       `json:"name,omitempty"`
	GivenName           string       `json:"given_name,omitempty"`
	FamilyName          string       `json:"family_name,omitempty"`
	MiddleName          string       `json:"middle_name,omitempty"`
	Nickname            string       `json:"nickname,omitempty"`
	PreferredUsername   string       `json:"preferred_username,omitempty"`
	Profile             string       `json:"profile,omitempty"`
	Picture             string       `json:"picture,omitempty"`
	Website             string       `json:"website,omitempty"`
	Email               string       `json:"email,omitempty"`
	EmailVerified       bool         `json:"email_verified,omitempty"`
	Gender              string       `json:"gender,omitempty"`
	BirthDate           string       `json:"birth_date,omitempty"`
	ZoneInfo            string       `json:"zone_info,omitempty"`
	Locale              string       `json:"locale,omitempty"`
	PhoneNumber         string       `json:"phone_number,omitempty"`
	PhoneNumberVerified bool         `json:"phone_number_verified,omitempty"`
	Address             *UserAddress `json:"address,omitempty"`
}

func (um *UserMetadata) Value() (marshal driver.Value, err error) {
	marshal, err = json.Marshal(um)
	return
}

func (um *UserMetadata) Scan(src interface{}) error {
	if bytes, ok := src.([]byte); ok {
		return json.Unmarshal(bytes, um)
	} else {
		return errors.New("failed to unmarshal value for UserMetadata")
	}
}
