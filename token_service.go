package core

import (
	"context"
	"database/sql"
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
	txn := ts.Db.WithContext(ctx)
	token := &models.TokensModel{
		RequestID:      reqId,
		ACSignature:    convertToNullString(signatures.GetACSignature()),
		ATSignature:    convertToNullString(signatures.GetATSignature()),
		RTSignature:    convertToNullString(signatures.GetRTSignature()),
		RTExpiry:       convertToNullTime(signatures.GetRTExpiry()),
		ATExpiry:       convertToNullTime(signatures.GetATExpiry()),
		ACExpiry:       convertToNullTime(signatures.GetACExpiry()),
		RequestProfile: &models.SavedProfile{Attributes: profile},
	}
	result := txn.Save(token)
	if result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func convertToNullTime(expiry time.Time) sql.NullTime {
	if expiry.IsZero() {
		return sql.NullTime{Valid: false}
	} else {
		return sql.NullTime{Valid: true, Time: expiry}
	}
}

func convertToNullString(signature string) sql.NullString {
	if signature == "" {
		return sql.NullString{Valid: false}
	} else {
		return sql.NullString{Valid: true, String: signature}
	}
}

func (ts *TokenStoreServiceImpl) GetProfileWithAuthCodeSign(ctx context.Context, signature string) (oidcsdk.RequestProfile, string, error) {
	txn := ts.Db.WithContext(ctx)
	token := &models.TokensModel{}
	result := txn.Find(token, "ac_signature = ?", signature)
	if result.Error != nil {
		return nil, "", result.Error
	}
	if result.RowsAffected != 1 {
		return nil, "", fmt.Errorf("authorization code not found")
	}
	if token.ACExpiry.Valid && token.ACExpiry.Time.Before(time.Now()) {
		return nil, "", fmt.Errorf("authorization code expired")
	}
	return token.RequestProfile.Attributes, token.RequestID, nil
}

func (ts *TokenStoreServiceImpl) GetProfileWithAccessTokenSign(ctx context.Context, signature string) (oidcsdk.RequestProfile, string, error) {
	txn := ts.Db.WithContext(ctx)
	token := &models.TokensModel{}
	result := txn.Find(token, "at_signature = ?", signature)
	if result.Error != nil {
		return nil, "", result.Error
	}
	if result.RowsAffected != 1 {
		return nil, "", fmt.Errorf("access token not found")
	}
	if token.ATExpiry.Valid && token.ATExpiry.Time.Before(time.Now()) {
		return nil, "", fmt.Errorf("access token expired")
	}
	return token.RequestProfile.Attributes, token.RequestID, nil
}

func (ts *TokenStoreServiceImpl) GetProfileWithRefreshTokenSign(ctx context.Context, signature string) (oidcsdk.RequestProfile, string, error) {
	txn := ts.Db.WithContext(ctx)
	token := &models.TokensModel{}
	result := txn.Find(token, "rt_signature = ?", signature)
	if result.Error != nil {
		return nil, "", result.Error
	}
	if result.RowsAffected != 1 {
		return nil, "", fmt.Errorf("refresh token not found")
	}
	if token.RTExpiry.Valid && token.RTExpiry.Time.Before(time.Now()) {
		return nil, "", fmt.Errorf("refresh token expired")
	}
	return token.RequestProfile.Attributes, token.RequestID, nil
}

func (ts *TokenStoreServiceImpl) InvalidateWithRequestID(ctx context.Context, reqID string, what uint8) (err error) {
	txn := ts.Db.WithContext(ctx)
	token := &models.TokensModel{}
	result := txn.Find(token, "request_id = ?", reqID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return fmt.Errorf("token not found with request id %s", reqID)
	}
	if token.RequestID != "" {
		if what&oidcsdk.ExpireRefreshToken > 0 {
			token.RTExpiry = sql.NullTime{Valid: true, Time: time.Now().Add(-10)}
		}
		if what&oidcsdk.ExpireAccessToken > 0 {
			token.ATExpiry = sql.NullTime{Valid: true, Time: time.Now().Add(-10)}
		}
		if what&oidcsdk.ExpireAuthorizationCode > 0 {
			token.ACExpiry = sql.NullTime{Valid: true, Time: time.Now().Add(-10)}
		}
	}
	return txn.Save(token).Error
}
