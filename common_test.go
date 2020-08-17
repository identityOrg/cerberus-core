package core

import (
	"github.com/jinzhu/gorm"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

var testDb *gorm.DB
var user *UserModel
var noCredUser *UserModel
var key *otp.Key

func init() {
	var err error
	testDb, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}
	testDb = testDb.Debug()
	testDb.AutoMigrate(&UserModel{}, &UserCredentials{}, &ServiceProviderModel{}, &ScopeModel{}, &ClaimModel{})
	err = testDb.Delete(&UserCredentials{}, "user_id = ?", 1).Error
	if err != nil {
		panic(err)
	}
	password, err := bcrypt.GenerateFromPassword([]byte("password"), 13)
	if err != nil {
		panic(err)
	}
	key, _ = totp.Generate(totp.GenerateOpts{
		Issuer:      "https://localhost:8080",
		AccountName: "user",
	})
	user = &UserModel{
		Username:     "user",
		EmailAddress: "user@domain.com",
		Metadata: &UserMetadata{
			Name: "Name Name",
		},
		Credentials: []UserCredentials{
			{
				Type:                CredTypePassword,
				Value:               string(password),
				InvalidAttemptCount: 2,
				Bocked:              false,
			},
			{
				Type:                CredTypeTOTP,
				Value:               key.Secret(),
				InvalidAttemptCount: 2,
				Bocked:              false,
			},
		},
		Inactive: false,
	}
	user.ID = 1
	err = testDb.Save(user).Error
	if err != nil {
		panic(err)
	}
	noCredUser = &UserModel{
		Username:     "nocred",
		EmailAddress: "nocred@domain.com",
		Metadata: &UserMetadata{
			Name: "No Cred",
		},
		Inactive: false,
	}
	noCredUser.ID = 2
	if testDb.Save(noCredUser).Error != nil {
		panic("noCred user not created")
	}
}
