package models

import "time"

type SecretChannelModel struct {
	BaseModel
	Name        string         `gorm:"column:name;unique_index:idx_channel_name" json:"name"`
	Algorithm   string         `gorm:"column:algorithm;unique_index:idx_alg_use" json:"algorithm"`
	Use         string         `gorm:"column:use;unique_index:idx_alg_use" json:"use"`
	ValidityDay uint           `gorm:"column:validity_day" json:"validity_day"`
	Secrets     []*SecretModel `gorm:"foreignKey:ChannelId" json:"secrets"`
}

func (sp SecretChannelModel) TableName() string {
	return "t_secret_channel"
}

type SecretModel struct {
	BaseModel
	KeyId     string    `gorm:"column:key_id" json:"key_id"`
	IssuedAt  time.Time `gorm:"column:issued_at" json:"issued_at"`
	ExpiresAt time.Time `gorm:"column:expires_at" json:"expires_at"`
	Value     []byte    `gorm:"column:value" json:"-"`
	ChannelId uint      `gorm:"column:channel_id" json:"channel_id"`
	Algorithm string    `gorm:"column:algorithm" json:"-"`
	Use       string    `gorm:"column:use" json:"-"`
}

func (sp SecretModel) TableName() string {
	return "t_secret"
}
