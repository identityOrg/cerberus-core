package core

import (
	"github.com/identityOrg/cerberus-core/models"
	"image"
)

type (
	IUserStoreService interface {
		IUserQueryService
		IUserCredentialsService
		IUserChangeService
		IUserCommonService
		IUserOTPService
	}
	IUserQueryService interface {
		GetUser(id uint) (user *models.UserModel, err error)
		FindUserByUsername(username string) (*models.UserModel, error)
		FindUserByEmail(email string) (*models.UserModel, error)
		FindAllUser(page uint, pageSize uint) ([]models.UserModel, uint, error)
	}
	IUserCredentialsService interface {
		SetPassword(id uint, password string) error
		GenerateTOTP(id uint, issuer string) (img image.Image, secret string, err error)
		ValidatePassword(id uint, password string) (err error)
		ValidateTOTP(id uint, code string) (err error)
	}
	IUserChangeService interface {
		ActivateUser(id uint) error
		DeactivateUser(id uint) error
		UsernameAvailable(username string) (available bool)
		ChangeUsername(id uint, username string) (err error)
		InitiateEmailChange(id uint, email string) (code string, err error)
		CompleteEmailChange(id uint, code string) (err error)
	}
	IUserCommonService interface {
		CreateUser(username string, email string, metadata *models.UserMetadata) (id uint, err error)
		UpdateUser(id uint, metadata *models.UserMetadata) (err error)
		PatchUser(id uint, metadata *models.UserMetadata) (err error)
		DeleteUser(id uint) (err error)
	}
	IUserOTPService interface {
		GenerateUserOTP(id uint, length uint8) (code string, err error)
		ValidateOTP(id uint, code string) (err error)
	}
)

const (
	CredTypePassword = 1
	CredTypeTOTP     = 2
)
