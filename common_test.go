package core

import (
	"github.com/jinzhu/gorm"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

var testDb *gorm.DB
var user *UserModel
var key *otp.Key

func init() {
	var err error
	testDb, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}
	testDb.AutoMigrate(&UserModel{}, &UserCredentials{}, &ServiceProviderModel{}, &ScopeModel{}, &ClaimModel{})
	password, err := bcrypt.GenerateFromPassword([]byte("password"), 13)
	key, _ = totp.Generate(totp.GenerateOpts{
		Issuer:      "https://localhost:8080",
		AccountName: "user",
	})
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
				Type:                CredTypePassword,
				Value:               string(password),
				FirstInvalidAttempt: nil,
				InvalidAttemptCount: 2,
				Bocked:              false,
			},
			{
				Type:                CredTypeTOTP,
				Value:               key.Secret(),
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
