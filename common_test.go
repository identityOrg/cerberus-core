package cerberus_models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

var testDb *gorm.DB

var user *UserModel

func init() {
	var err error
	testDb, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}
	testDb.AutoMigrate(&UserModel{}, &UserCredentials{}, &ServiceProviderModel{}, &ScopeModel{}, &ClaimModel{})
	password, err := bcrypt.GenerateFromPassword([]byte("password"), 13)
	user = &UserModel{
		BaseModel: BaseModel{
			ID: 1,
		},
		Username:     "user",
		EmailAddress: "user@domain.com",
		Metadata: &UserMetadata{
			Name: "Name Name",
		},
		Credentials: []UserCredentials{
			{
				Type:                0,
				Value:               string(password),
				FirstInvalidAttempt: nil,
				InvalidAttemptCount: 2,
				Bocked:              false,
			},
		},
		Inactive: false,
	}
	err = testDb.Save(user).Error
	if err != nil {
		panic(err)
	}
}
