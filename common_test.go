package core

import (
	"github.com/identityOrg/cerberus-core/models"
	"github.com/jinzhu/gorm"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

var TestDb *gorm.DB
var TestUser *models.UserModel
var TestNoCredUser *models.UserModel
var TestNoCredUser2 *models.UserModel
var key *otp.Key

func init() {
	var err error
	TestDb, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}
	TestDb = TestDb.Debug()
	TestDb.AutoMigrate(&models.UserModel{}, &models.UserCredentials{}, &models.ServiceProviderModel{}, &models.ScopeModel{}, &models.ClaimModel{})
	err = TestDb.Delete(&models.UserCredentials{}).Error
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
	TestUser = &models.UserModel{
		Username:     "user",
		EmailAddress: "user@domain.com",
		Metadata: &models.UserMetadata{
			Name: "Name Name",
		},
		Credentials: []models.UserCredentials{
			{
				Type:   CredTypePassword,
				Value:  string(password),
				Bocked: false,
			},
			{
				Type:   CredTypeTOTP,
				Value:  key.Secret(),
				Bocked: false,
			},
		},
		Inactive: false,
	}
	TestUser.ID = 1
	err = TestDb.Save(TestUser).Error
	if err != nil {
		panic(err)
	}
	TestNoCredUser = &models.UserModel{
		Username:     "nocred",
		EmailAddress: "nocred@domain.com",
		Metadata: &models.UserMetadata{
			Name: "No Cred",
		},
		Inactive: false,
	}
	TestNoCredUser.ID = 2
	if TestDb.Save(TestNoCredUser).Error != nil {
		panic("noCred user not created")
	}
	TestNoCredUser2 = &models.UserModel{
		Username:     "nocred2",
		EmailAddress: "nocred2@domain.com",
		Metadata: &models.UserMetadata{
			Name: "No Cred2",
		},
		Inactive: false,
	}
	TestNoCredUser2.ID = 3
	if TestDb.Save(TestNoCredUser2).Error != nil {
		panic("noCred user not created")
	}
}
