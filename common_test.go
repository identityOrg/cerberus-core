package core

import (
	"github.com/google/uuid"
	"github.com/identityOrg/cerberus-core/models"
	"github.com/jinzhu/gorm"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

var (
	TestDb          *gorm.DB
	TestUser        *models.UserModel
	TestNoCredUser  *models.UserModel
	TestNoCredUser2 *models.UserModel
	key             *otp.Key
	TestSP          *models.ServiceProviderModel
	TestSP2         *models.ServiceProviderModel
)

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
	createTestUser()
	createNoCredUser()
	createNoCredUser2()
	createTestSP()
}

func createTestSP() {
	TestSP = &models.ServiceProviderModel{
		Name:         "TestSP",
		Description:  "A Service Provider for test",
		ClientID:     uuid.New().String(),
		ClientSecret: uuid.New().String(),
		Active:       true,
		Metadata: &models.ServiceProviderMetadata{
			RedirectUris:    []string{"http://localhost:9090/redirect"},
			ResponseTypes:   []string{"code", "token", "token id_token", "id_token", "code id_token"},
			GrantTypes:      []string{"authorization_code", "password"},
			ApplicationType: "web",
		},
	}
	TestSP.ID = 1
	err := TestDb.Save(TestSP).Error
	if err != nil {
		panic(err)
	}
}

func createNoCredUser2() {
	TestNoCredUser2 = &models.UserModel{
		Username:     "nocred2",
		EmailAddress: "nocred2@domain.com",
		Metadata: &models.UserMetadata{
			Name: "No Cred2",
		},
		Inactive: false,
	}
	TestNoCredUser2.ID = 3
	err := TestDb.Save(TestNoCredUser2).Error
	if err != nil {
		panic("noCred user not created:" + err.Error())
	}
}

func createNoCredUser() {
	TestNoCredUser = &models.UserModel{
		Username:     "nocred",
		EmailAddress: "nocred@domain.com",
		Metadata: &models.UserMetadata{
			Name: "No Cred",
		},
		Inactive: false,
	}
	TestNoCredUser.ID = 2
	err := TestDb.Save(TestNoCredUser).Error
	if err != nil {
		panic("noCred user not created:" + err.Error())
	}
}

func createTestUser() {
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
}
