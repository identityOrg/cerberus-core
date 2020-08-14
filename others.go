package core

type ScopeModel struct {
	BaseModel
	Name        string `gorm:"column:name;size:256" json:"name,omitempty"`
	Description string `gorm:"column:description;size:1024" json:"description,omitempty"`
}

func (sm ScopeModel) TableName() string {
	return "t_scopes"
}

type ClaimModel struct {
	BaseModel
	Name        string `gorm:"column:name;size:256" json:"name,omitempty"`
	Description string `gorm:"column:description;size:1024" json:"description,omitempty"`
}

func (cm ClaimModel) TableName() string {
	return "t_claims"
}
