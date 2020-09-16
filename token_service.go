package core

import (
	"context"
	"fmt"
	"github.com/identityOrg/cerberus-core/models"
	"github.com/jinzhu/gorm"
	"time"
)

type TokenStoreServiceImpl struct {
	Db *gorm.DB
}

func (ts *TokenStoreServiceImpl) BeginTransaction(ctx context.Context, readOnly bool) context.Context {
	return beginTransaction(ctx, readOnly, ts.Db)
}

func (ts *TokenStoreServiceImpl) CommitTransaction(ctx context.Context) context.Context {
	return commitTransaction(ctx)
}

func (ts *TokenStoreServiceImpl) RollbackTransaction(ctx context.Context) context.Context {
	return rollbackTransaction(ctx)
}

func (ts *TokenStoreServiceImpl) StoreTokenProfile(ctx context.Context, reqId string, signatures ITokenSignatures, profile map[string]string) (err error) {
	txn := getTransaction(ctx)
	token := &models.TokensModel{
		RequestID:      reqId,
		ACSignature:    signatures.GetACSignature(),
		ATSignature:    signatures.GetATSignature(),
		RTSignature:    signatures.GetRTSignature(),
		RTExpiry:       signatures.GetRTExpiry(),
		ATExpiry:       signatures.GetATExpiry(),
		ACExpiry:       signatures.GetACExpiry(),
		RequestProfile: profile,
	}
	result := txn.Save(token)
	if result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (ts *TokenStoreServiceImpl) GetProfileWithAuthCodeSign(ctx context.Context, signature string) (profile map[string]string, reqId string, err error) {
	txn := getTransaction(ctx)
	token := &models.TokensModel{}
	result := txn.Find(token, "ac_signature = ?", signature)
	if result.Error != nil {
		return nil, "", result.Error
	}
	if result.RecordNotFound() {
		return nil, "", fmt.Errorf("authorization code not found")
	}
	if token.ACExpiry.After(time.Now()) {
		return nil, "", fmt.Errorf("authorization code expired")
	}
	return token.RequestProfile, token.RequestID, nil
}

func (ts *TokenStoreServiceImpl) GetProfileWithAccessTokenSign(ctx context.Context, signature string) (profile map[string]string, reqId string, err error) {
	txn := getTransaction(ctx)
	token := &models.TokensModel{}
	result := txn.Find(token, "at_signature = ?", signature)
	if result.Error != nil {
		return nil, "", result.Error
	}
	if result.RecordNotFound() {
		return nil, "", fmt.Errorf("access token not found")
	}
	if token.ACExpiry.After(time.Now()) {
		return nil, "", fmt.Errorf("access token expired")
	}
	return token.RequestProfile, token.RequestID, nil
}

func (ts *TokenStoreServiceImpl) GetProfileWithRefreshTokenSign(ctx context.Context, signature string) (profile map[string]string, reqId string, err error) {
	txn := getTransaction(ctx)
	token := &models.TokensModel{}
	result := txn.Find(token, "rt_signature = ?", signature)
	if result.Error != nil {
		return nil, "", result.Error
	}
	if result.RecordNotFound() {
		return nil, "", fmt.Errorf("refresh token not found")
	}
	if token.ACExpiry.After(time.Now()) {
		return nil, "", fmt.Errorf("refresh token expired")
	}
	return token.RequestProfile, token.RequestID, nil
}

func (ts *TokenStoreServiceImpl) InvalidateWithRequestID(ctx context.Context, reqID string, what uint8) (err error) {
	txn := getTransaction(ctx)
	token := &models.TokensModel{}
	result := txn.Find(token, "request_id = ?", reqID)
	if result.Error != nil {
		return result.Error
	}
	if token.RequestID != "" {
		if what&expireRefreshToken > 0 {
			token.RTExpiry = time.Now()
		}
		if what&expireAccessToken > 0 {
			token.ATExpiry = time.Now()
		}
		if what&expireAuthorizationCode > 0 {
			token.ACExpiry = time.Now()
		}
	}
	return txn.Save(token).Error
}

const (
	expireAuthorizationCode = 1
	expireAccessToken       = 2
	expireRefreshToken      = 4
)
