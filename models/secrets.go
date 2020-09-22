package models

import "time"

type SecretChannelModel struct {
	BaseModel
	Name        string         `sql:"column:name" json:"name"`
	Algorithm   string         `sql:"column:algorithm" json:"algorithm"`
	Use         string         `sql:"column:use" json:"use"`
	ValidityDay uint           `sql:"column:validity_day;unique_index" json:"validity_day"`
	Secrets     []*SecretModel `sql:"foreignKey:ChannelId" json:"secrets"`
}

func (sp SecretChannelModel) TableName() string {
	return "t_secret_channel"
}

type SecretModel struct {
	BaseModel
	KeyId     string    `sql:"column:key_id" json:"key_id"`
	IssuedAt  time.Time `sql:"column:issued_at" json:"issued_at"`
	ExpiresAt time.Time `sql:"column:expires_at" json:"expires_at"`
	Value     []byte    `sql:"column:value;type:lob" json:"-"`
	ChannelId uint      `sql:"column:channel_id" json:"channel_id"`
	Algorithm string    `sql:"column:algorithm" json:"-"`
	Use       string    `sql:"column:use" json:"-"`
}

func (sp SecretModel) TableName() string {
	return "t_secret"
}
