package core

import (
	"github.com/identityOrg/cerberus-core/models"
	"gorm.io/gorm"
)

var StandardScopes = map[string]string{
	"openid":  "OpenID Scope",
	"profile": "Profile Scope",
	"email":   "Email Scope",
	"phone":   "Phone Scope",
	"address": "Address Scope",
}
var StandardClaims = map[string][]struct {
	Name string
	Desc string
}{
	"openid": {
		{Name: "sub", Desc: "Subject"},
	},
	"profile": {
		{Name: "name", Desc: "Name"},
		{Name: "given_name", Desc: "Given Name"},
		{Name: "family_name", Desc: "Family Name"},
		{Name: "middle_name", Desc: "Middle Name"},
		{Name: "nickname", Desc: "Nickname"},
		{Name: "preferred_username", Desc: "Preferred Username"},
		{Name: "profile", Desc: "Profile"},
		{Name: "picture", Desc: "Picture"},
		{Name: "website", Desc: "Website"},
		{Name: "gender", Desc: "Gender"},
		{Name: "birthdate", Desc: "Birth Date"},
		{Name: "zoneinfo", Desc: "Zone Info"},
		{Name: "locale", Desc: "Locale"},
	},
	"email": {
		{Name: "email", Desc: "Email"},
		{Name: "email_verified", Desc: "Email Verified"},
	},
	"phone": {
		{Name: "phone_number", Desc: "Phone Number"},
		{Name: "phone_number_verified", Desc: "Phone Number Verified"},
	},
	"address": {
		{Name: "address", Desc: "Address"},
	},
}

func InitializeDefaultScope(db *gorm.DB) error {
	for name, desc := range StandardScopes {
		scope := &models.ScopeModel{
			Name:        name,
			Description: desc,
		}
		saveResult := db.FirstOrCreate(scope, "name = ?", name)
		if saveResult.Error != nil {
			return saveResult.Error
		}
		ass := db.Model(scope).Association("Claims")
		if mapping, ok := StandardClaims[name]; ok {
			for _, mm := range mapping {
				claim := &models.ClaimModel{
					Name:        mm.Name,
					Description: mm.Desc,
				}
				err := db.FirstOrCreate(claim, "name = ?", mm.Name).Error
				if err != nil {
					return err
				}
				err = ass.Append(claim)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
