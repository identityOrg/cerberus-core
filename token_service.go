package core

import (
	"context"
	"fmt"
	"github.com/identityOrg/cerberus-core/models"
	"github.com/identityOrg/oidcsdk"
	"gorm.io/gorm"
	"time"
)

type TokenStoreServiceImpl struct {
	Db *gorm.DB
}

func NewTokenStoreServiceImpl(db *gorm.DB) *TokenStoreServiceImpl {
	return &TokenStoreServiceImpl{Db: db}
}

func (ts *TokenStoreServiceImpl) StoreTokenProfile(ctx context.Context, reqId string, signatures oidcsdk.ITokenSignatures, profile oidcsdk.RequestProfile) (err error) {
	txn := ts.Db
	token := &models.TokensModel{
		RequestID:      reqId,
		ACSignature:    signatures.GetACSignature(),
		ATSignature:    signatures.GetATSignature(),
		RTSignature:    signatures.GetRTSignature(),
		RTExpiry:       signatures.GetRTExpiry(),
		ATExpiry:       signatures.GetATExpiry(),
		ACExpiry:       signatures.GetACExpiry(),
		RequestProfile: &models.SavedProfile{Attributes: profile},
	}
	result := txn.Save(token)
	if result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (ts *TokenStoreServiceImpl) GetProfileWithAuthCodeSign(ctx context.Context, signature string) (oidcsdk.RequestProfile, string, error) {
	txn := ts.Db
	token := &models.TokensModel{}
	result := txn.Find(token, "ac_signature = ?", signature)
	if result.Error != nil {
		return nil, "", result.Error
	}
	if result.RowsAffected != 1 {
		return nil, "", fmt.Errorf("authorization code not found")
	}
	if token.ACExpiry.Before(time.Now()) {
		return nil, "", fmt.Errorf("authorization code expired")
	}
	return token.RequestProfile.Attributes, token.RequestID, nil
}

func (ts *TokenStoreServiceImpl) GetProfileWithAccessTokenSign(ctx context.Context, signature string) (oidcsdk.RequestProfile, string, error) {
	txn := ts.Db
	token := &models.TokensModel{}
	result := txn.Find(token, "at_signature = ?", signature)
	if result.Error != nil {
		return nil, "", result.Error
	}
	if result.Error != nil {
		return nil, "", result.Error
	}
	if token.ATExpiry.Before(time.Now()) {
		return nil, "", fmt.Errorf("access token expired")
	}
	return token.RequestProfile.Attributes, token.RequestID, nil
}

func (ts *TokenStoreServiceImpl) GetProfileWithRefreshTokenSign(ctx context.Context, signature string) (oidcsdk.RequestProfile, string, error) {
	txn := ts.Db
	token := &models.TokensModel{}
	result := txn.Find(token, "rt_signature = ?", signature)
	if result.Error != nil {
		return nil, "", result.Error
	}
	if result.Error != nil {
		return nil, "", result.Error
	}
	if token.RTExpiry.Before(time.Now()) {
		return nil, "", fmt.Errorf("refresh token expired")
	}
	return token.RequestProfile.Attributes, token.RequestID, nil
}

func (ts *TokenStoreServiceImpl) InvalidateWithRequestID(ctx context.Context, reqID string, what uint8) (err error) {
	txn := ts.Db
	token := &models.TokensModel{}
	result := txn.Find(token, "request_id = ?", reqID)
	if result.Error != nil {
		return result.Error
	}
	if token.RequestID != "" {
		if what&oidcsdk.ExpireRefreshToken > 0 {
			token.RTExpiry = time.Now().Add(-10)
		}
		if what&oidcsdk.ExpireAccessToken > 0 {
			token.ATExpiry = time.Now().Add(-10)
		}
		if what&oidcsdk.ExpireAuthorizationCode > 0 {
			token.ACExpiry = time.Now().Add(-10)
		}
	}
	return txn.Save(token).Error
}
