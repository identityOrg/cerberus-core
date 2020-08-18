package core

import (
	"github.com/identityOrg/cerberus-core/models"
	"image"
)

type IUserStoreService interface {
	FindUserByUsername(username string) (*models.UserModel, error)
	FindUserByEmail(email string) (*models.UserModel, error)
	FindAllUser(page uint, pageSize uint) ([]models.UserModel, uint, error)
	ActivateUser(id uint) error
	DeactivateUser(id uint) error
	SetPassword(id uint, password string) error
	GenerateTOTP(id uint, issuer string) (img image.Image, secret string, err error)
	ValidatePassword(id uint, password string) (err error)
	ValidateTOTP(id uint, code string) (err error)
}

const (
	CredTypePassword = 1
	CredTypeTOTP     = 2
)
