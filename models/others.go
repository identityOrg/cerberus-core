package models

type ScopeModel struct {
	BaseModel
	Name        string        `gorm:"column:name;size:256;index:idx_scope_name,unique" json:"name,omitempty"`
	Description string        `gorm:"column:description;size:1024" json:"description,omitempty"`
	Claims      []*ClaimModel `gorm:"many2many:t_scope_claim"`
}

func (sm ScopeModel) TableName() string {
	return "t_scopes"
}

type ClaimModel struct {
	BaseModel
	Name        string `gorm:"column:name;size:256;index:idx_claim_name,unique" json:"name,omitempty"`
	Description string `gorm:"column:description;size:1024" json:"description,omitempty"`
}

func (cm ClaimModel) TableName() string {
	return "t_claims"
}
