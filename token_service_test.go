package core

import (
	"context"
	"github.com/google/uuid"
	"github.com/identityOrg/oidcsdk"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTokenStoreServiceImpl(t *testing.T) {
	var tokenService = NewTokenStoreServiceImpl(TestDb)
	tokenService.Db = beginTransaction(context.Background(), tokenService.Db)
	ctx := context.Background()
	t.Run("ensure store", func(t *testing.T) {
		signMock := NewTokenSignMock(time.Now().Add(time.Minute * 10))
		profile := make(map[string]string)
		profile["key"] = "value"
		reqId := uuid.New().String()
		err := tokenService.StoreTokenProfile(ctx, reqId, signMock, profile)
		if assert.NoError(t, err) {
			t.Run("find by AT", func(t *testing.T) {
				profile, id, err := tokenService.GetProfileWithAccessTokenSign(ctx, signMock.GetATSignature())
				if assert.NoError(t, err) {
					if assert.Equal(t, reqId, id) {
						assert.Equal(t, "value", profile["key"])
					}
				}
			})
			t.Run("find by RT", func(t *testing.T) {
				profile, id, err := tokenService.GetProfileWithRefreshTokenSign(ctx, signMock.GetRTSignature())
				if assert.NoError(t, err) {
					if assert.Equal(t, reqId, id) {
						assert.Equal(t, "value", profile["key"])
					}
				}
			})
			t.Run("find by AC", func(t *testing.T) {
				profile, id, err := tokenService.GetProfileWithAuthCodeSign(ctx, signMock.GetACSignature())
				if assert.NoError(t, err) {
					if assert.Equal(t, reqId, id) {
						assert.Equal(t, "value", profile["key"])
					}
				}
			})
			t.Run("invalidate AC", func(t *testing.T) {
				err = tokenService.InvalidateWithRequestID(ctx, reqId, oidcsdk.ExpireAuthorizationCode)
				if assert.NoError(t, err) {
					_, _, err := tokenService.GetProfileWithAuthCodeSign(ctx, signMock.GetACSignature())
					if assert.Error(t, err) {
						assert.EqualError(t, err, "authorization code expired")
					}
					_, _, err = tokenService.GetProfileWithRefreshTokenSign(ctx, signMock.GetRTSignature())
					assert.NoError(t, err)
				}
			})
		}
	})
	t.Run("negative test", func(t *testing.T) {
		signMock := NewTokenSignMock(time.Now().Add(-10))
		profile := make(map[string]string)
		profile["key"] = "value"
		reqId := uuid.New().String()
		err := tokenService.StoreTokenProfile(ctx, reqId, signMock, profile)
		if assert.NoError(t, err) {
			t.Run("expired", func(t *testing.T) {
				_, _, err := tokenService.GetProfileWithAuthCodeSign(ctx, signMock.GetACSignature())
				if assert.Error(t, err) {
					assert.EqualError(t, err, "authorization code expired")
				}
			})
			t.Run("non existing", func(t *testing.T) {
				_, _, err := tokenService.GetProfileWithAuthCodeSign(ctx, "not existing code")
				if assert.Error(t, err) {
					assert.EqualError(t, err, "authorization code not found")
				}
			})
		}
	})
	rollbackTransaction(tokenService.Db)
}

type TokenSignMock struct {
	Expiry time.Time
	ATSign string
	RTSign string
	ACSign string
}

func NewTokenSignMock(expiry time.Time) *TokenSignMock {
	return &TokenSignMock{
		Expiry: expiry,
		ATSign: uuid.New().String(),
		RTSign: uuid.New().String(),
		ACSign: uuid.New().String(),
	}
}

func (ts TokenSignMock) GetACSignature() string {
	return ts.ACSign
}

func (ts TokenSignMock) GetATSignature() string {
	return ts.ATSign
}

func (ts TokenSignMock) GetRTSignature() string {
	return ts.RTSign
}

func (ts TokenSignMock) GetACExpiry() time.Time {
	return ts.Expiry
}

func (ts TokenSignMock) GetATExpiry() time.Time {
	return ts.Expiry
}

func (ts TokenSignMock) GetRTExpiry() time.Time {
	return ts.Expiry
}
