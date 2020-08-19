package core

import (
	"context"
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
		ITransactionalStore
	}
	IUserQueryService interface {
		GetUser(ctx context.Context, id uint) (user *models.UserModel, err error)
		FindUserByUsername(ctx context.Context, username string) (*models.UserModel, error)
		FindUserByEmail(ctx context.Context, email string) (*models.UserModel, error)
		FindAllUser(ctx context.Context, page uint, pageSize uint) ([]models.UserModel, uint, error)
	}
	IUserCredentialsService interface {
		SetPassword(ctx context.Context, id uint, password string) error
		GenerateTOTP(ctx context.Context, id uint, issuer string) (img image.Image, secret string, err error)
		ValidatePassword(ctx context.Context, id uint, password string) (err error)
		ValidateTOTP(ctx context.Context, id uint, code string) (err error)
	}
	IUserChangeService interface {
		ActivateUser(ctx context.Context, id uint) error
		DeactivateUser(ctx context.Context, id uint) error
		UsernameAvailable(ctx context.Context, username string) (available bool)
		ChangeUsername(ctx context.Context, id uint, username string) (err error)
		InitiateEmailChange(ctx context.Context, id uint, email string) (code string, err error)
		CompleteEmailChange(ctx context.Context, id uint, code string) (err error)
	}
	IUserCommonService interface {
		CreateUser(ctx context.Context, username string, email string, metadata *models.UserMetadata) (id uint, err error)
		UpdateUser(ctx context.Context, id uint, metadata *models.UserMetadata) (err error)
		PatchUser(ctx context.Context, id uint, metadata *models.UserMetadata) (err error)
		DeleteUser(ctx context.Context, id uint) (err error)
	}
	IUserOTPService interface {
		GenerateUserOTP(ctx context.Context, id uint, length uint8) (code string, err error)
		ValidateOTP(ctx context.Context, id uint, code string) (err error)
	}
	ITransactionalStore interface {
		BeginTransaction(ctx context.Context, readOnly bool) context.Context
		CommitTransaction(ctx context.Context) context.Context
		RollbackTransaction(ctx context.Context) context.Context
	}
)

const (
	CredTypePassword = 1
	CredTypeTOTP     = 2
)
