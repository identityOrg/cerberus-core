package core

import (
	"context"
	"github.com/identityOrg/cerberus-core/models"
	"github.com/identityOrg/oidcsdk"
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
	ISPStoreService interface {
		ISPCommonService
		ISPUpdateService
		ISPCredentialService
		ISPQueryService
	}
	ISPCommonService interface {
		CreateSP(ctx context.Context, clientName string, description string, metadata *models.ServiceProviderMetadata) (id uint, err error)
		UpdateSP(ctx context.Context, id uint, metadata *models.ServiceProviderMetadata) (err error)
		PatchSP(ctx context.Context, id uint, metadata *models.ServiceProviderMetadata) (err error)
		DeleteSP(ctx context.Context, id uint) (err error)
	}
	ISPUpdateService interface {
		ActivateSP(ctx context.Context, id uint) error
		DeactivateSP(ctx context.Context, id uint) error
	}
	ISPCredentialService interface {
		ResetClientCredentials(ctx context.Context, id uint) (clientId, clientSecret string, err error)
		ValidateClientCredentials(ctx context.Context, clientId, clientSecret string) (id uint, err error)
		ValidateSecretSignature(ctx context.Context, token string) (id uint, err error)
		ValidatePrivateKeySignature(ctx context.Context, token string) (id uint, err error)
	}
	ISPQueryService interface {
		GetSP(ctx context.Context, id uint) (sp *models.ServiceProviderModel, err error)
		FindSPByClientId(ctx context.Context, clientId string) (sp *models.ServiceProviderModel, err error)
		FindSPByName(ctx context.Context, name string) (sp *models.ServiceProviderModel, err error)
		FindAllSP(ctx context.Context, page uint, pageSize uint) (sps []models.ServiceProviderModel, count uint, err error)
	}
	ITextEncrypts interface {
		EncryptText(ctx context.Context, text string) (cypherText string, err error)
	}
	ITextDecrypts interface {
		DecryptText(ctx context.Context, cypherText string) (text string, err error)
	}
	ITokenStoreService interface {
		oidcsdk.ITokenStore
	}
	IScopeClaimStoreService interface {
		IScopeOperations
		IClaimOperations
	}
	IClaimOperations interface {
		CreateClaim(ctx context.Context, name string, description string) (id uint, err error)
		FindClaimByName(ctx context.Context, name string) (*models.ClaimModel, error)
		GetClaim(ctx context.Context, id uint) (*models.ClaimModel, error)
		GetAllClaims(ctx context.Context, page uint, pageSize uint) ([]*models.ClaimModel, uint, error)
		UpdateClaim(ctx context.Context, id uint, description string) error
		DeleteClaim(ctx context.Context, id uint) error
	}
	IScopeOperations interface {
		CreateScope(ctx context.Context, name string, description string) (id uint, err error)
		FindScopeByName(ctx context.Context, name string) (*models.ScopeModel, error)
		GetScope(ctx context.Context, id uint) (*models.ScopeModel, error)
		GetAllScopes(ctx context.Context, page uint, pageSize uint) ([]*models.ScopeModel, uint, error)
		UpdateScope(ctx context.Context, id uint, description string) error
		DeleteScope(ctx context.Context, id uint) error
		AddClaimToScope(ctx context.Context, scopeId uint, claimId uint) error
		RemoveClaimFromScope(ctx context.Context, scopeId uint, claimId uint) error
	}
	ISecretStoreService interface {
		oidcsdk.ISecretStore
		ISecretChannelManager
	}
	ISecretChannelManager interface {
		CreateChannel(ctx context.Context, name string, algorithm string, use string, validityDay uint) (uint, error)
		GetAllChannels(ctx context.Context) ([]*models.SecretChannelModel, error)
		GetChannel(ctx context.Context, channelId uint) (*models.SecretChannelModel, error)
		DeleteChannel(ctx context.Context, channelId uint) error
		RenewSecret(ctx context.Context, channelId uint) error
	}
)

const (
	CredTypePassword = 1
	CredTypeTOTP     = 2
)
