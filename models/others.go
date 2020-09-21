package models

type ScopeModel struct {
	BaseModel
	Name        string        `sql:"column:name;size:256" json:"name,omitempty"`
	Description string        `sql:"column:description;size:1024" json:"description,omitempty"`
	Claims      []*ClaimModel `sql:"many2many:t_scope_claim"`
}

func (sm ScopeModel) TableName() string {
	return "t_scopes"
}

type ClaimModel struct {
	BaseModel
	Name        string `sql:"column:name;size:256" json:"name,omitempty"`
	Description string `sql:"column:description;size:1024" json:"description,omitempty"`
}

func (cm ClaimModel) TableName() string {
	return "t_claims"
}
