package core

import (
	"context"
	"encoding/base32"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTOTP(t *testing.T) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "https://localhost:8080",
		AccountName: "username",
		Period:      30,
	})
	if err != nil {
		t.Error(err)
	}
	println(key.Secret())

	code, err := totp.GenerateCode(key.Secret(), time.Now())
	if err != nil {
		t.Error(err)
	}
	println(code)
	println(totp.Validate(code, key.Secret()))
}

var b32NoPadding = base32.StdEncoding.WithPadding(base32.NoPadding)

func TestRerenderPng(t *testing.T) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "https://localhost:8080",
		AccountName: "username",
		Period:      30,
	})
	if err != nil {
		t.Error(err)
	}
	println(key.Secret())

	decoded, _ := b32NoPadding.DecodeString(key.Secret())
	key2, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "https://localhost:8080",
		AccountName: "username",
		Secret:      decoded,
	})
	if err != nil {
		t.Error(err)
	}
	println(key2.Secret())
}

func TestUserStoreServiceImpl_ActivateDeactivateUser(t *testing.T) {
	ctx := context.Background()
	userStoreService := NewUserStoreService(TestDb, 3, 5*time.Minute)
	ctx = userStoreService.BeginTransaction(ctx, true)
	t.Run("de-activate", func(t *testing.T) {
		err := userStoreService.DeactivateUser(ctx, TestUser.ID)
		assert.NoError(t, err)
	})
	t.Run("activate", func(t *testing.T) {
		err := userStoreService.ActivateUser(ctx, TestUser.ID)
		assert.NoError(t, err)
	})
	t.Run("user not exists", func(t *testing.T) {
		err := userStoreService.ActivateUser(ctx, 1000)
		if assert.Error(t, err) {
			assert.NotNil(t, err)
		}
	})
	userStoreService.RollbackTransaction(ctx)
}

func TestUserStoreServiceImpl_FindAllUser(t *testing.T) {
	ctx := context.Background()
	userStoreService := NewUserStoreService(TestDb, 3, 5*time.Minute)
	ctx = userStoreService.BeginTransaction(ctx, true)
	allUser, count, err := userStoreService.FindAllUser(ctx, 0, 5)
	assert.Nil(t, err)
	if assert.Equal(t, uint(3), count) {
		assert.Equal(t, uint(1), allUser[0].ID)
	}
	userStoreService.RollbackTransaction(ctx)
}

func TestUserStoreServiceImpl_ValidatePassword(t *testing.T) {
	ctx := context.Background()
	userStoreService := NewUserStoreService(TestDb, 3, 5*time.Minute)
	ctx = userStoreService.BeginTransaction(ctx, true)
	t.Run("valid", func(t *testing.T) {
		err := userStoreService.ValidatePassword(ctx, 1, "password")
		assert.Nil(t, err)
	})
	t.Run("invalid", func(t *testing.T) {
		err := userStoreService.ValidatePassword(ctx, 1, "password1")
		assert.Error(t, err)
	})
	t.Run("user not exists", func(t *testing.T) {
		err := userStoreService.ValidatePassword(ctx, 2000, "password")
		assert.Error(t, err)
	})
	t.Run("cred not found", func(t *testing.T) {
		err := userStoreService.ValidatePassword(ctx, TestNoCredUser2.ID, "password")
		if assert.Error(t, err) {
			assert.EqualError(t, err, "credential not found")
		}
	})
	userStoreService.RollbackTransaction(ctx)
}

func TestUserStoreServiceImpl_FindUserByEmail(t *testing.T) {
	ctx := context.Background()
	userStoreService := NewUserStoreService(TestDb, 3, 5*time.Minute)
	ctx = userStoreService.BeginTransaction(ctx, true)
	t.Run("found", func(t *testing.T) {
		foundUser, err := userStoreService.FindUserByEmail(ctx, TestUser.EmailAddress)
		assert.NoError(t, err)
		assert.Equal(t, TestUser.ID, foundUser.ID)
	})
	t.Run("not found", func(t *testing.T) {
		foundUser, err := userStoreService.FindUserByEmail(ctx, "invalid@domain.com")
		assert.Error(t, err)
		assert.Nil(t, foundUser)
	})
	userStoreService.RollbackTransaction(ctx)
}

func TestUserStoreServiceImpl_FindUserByUsername(t *testing.T) {
	ctx := context.Background()
	userStoreService := NewUserStoreService(TestDb, 3, 5*time.Minute)
	ctx = userStoreService.BeginTransaction(ctx, true)
	t.Run("found", func(t *testing.T) {
		foundUser, err := userStoreService.FindUserByUsername(ctx, TestUser.Username)
		assert.NoError(t, err)
		assert.Equal(t, TestUser.ID, foundUser.ID)
	})
	t.Run("not found", func(t *testing.T) {
		foundUser, err := userStoreService.FindUserByUsername(ctx, "invalid")
		assert.Error(t, err)
		assert.Nil(t, foundUser)
	})
	userStoreService.RollbackTransaction(ctx)
}

func TestUserStoreServiceImpl_ValidateTOTP(t *testing.T) {
	ctx := context.Background()
	userStoreService := NewUserStoreService(TestDb, 3, 5*time.Minute)
	ctx = userStoreService.BeginTransaction(ctx, true)
	t.Run("valid", func(t *testing.T) {
		code, err := totp.GenerateCode(TestUser.Credentials[1].Value, time.Now())
		if assert.NoError(t, err) {
			err = userStoreService.ValidateTOTP(ctx, TestUser.ID, code)
			assert.NoError(t, err)
		}
	})
	t.Run("invalid", func(t *testing.T) {
		err := userStoreService.ValidateTOTP(ctx, TestUser.ID, "code")
		assert.Error(t, err)
	})
	t.Run("no user", func(t *testing.T) {
		err := userStoreService.ValidateTOTP(ctx, 2000, "code")
		assert.Error(t, err)
	})
	userStoreService.RollbackTransaction(ctx)
}

func TestUserStoreServiceImpl_SetPassword(t *testing.T) {
	ctx := context.Background()
	userStoreService := NewUserStoreService(TestDb, 3, 5*time.Minute)
	ctx = userStoreService.BeginTransaction(ctx, true)
	t.Run("success", func(t *testing.T) {
		err := userStoreService.SetPassword(ctx, TestNoCredUser.ID, "new password")
		assert.NoError(t, err)
	})
	t.Run("user not found", func(t *testing.T) {
		err := userStoreService.SetPassword(ctx, 2000, "new password")
		assert.Error(t, err)
	})
	userStoreService.RollbackTransaction(ctx)
}

func TestUserStoreServiceImpl_GenerateTOTP(t *testing.T) {
	ctx := context.Background()
	userStoreService := NewUserStoreService(TestDb, 3, 5*time.Minute)
	ctx = userStoreService.BeginTransaction(ctx, true)
	t.Run("success", func(t *testing.T) {
		image, secret, err := userStoreService.GenerateTOTP(ctx, TestNoCredUser.ID, "cerberus")
		if assert.NoError(t, err) {
			assert.NotNil(t, image)
			assert.NotEqual(t, "", secret)
		}
	})
	t.Run("user not found", func(t *testing.T) {
		image, secret, err := userStoreService.GenerateTOTP(ctx, 2000, "cerberus")
		if assert.Error(t, err) {
			assert.Nil(t, image)
			assert.EqualError(t, err, "user not found")
			assert.Equal(t, "", secret)
		}
	})
	userStoreService.RollbackTransaction(ctx)
}

func TestUserStoreServiceImpl_Credential_Block_Unblock(t *testing.T) {
	ctx := context.Background()
	userStoreService := NewUserStoreService(TestDb, 1, 5*time.Minute)
	ctx = userStoreService.BeginTransaction(ctx, true)
	t.Run("blocked", func(t *testing.T) {
		err := userStoreService.SetPassword(ctx, TestUser.ID, "other password")
		if assert.NoError(t, err) {
			err = userStoreService.ValidatePassword(ctx, TestUser.ID, "invalid")
			if assert.Error(t, err) {
				assert.EqualError(t, err, "password mismatch")
			}
		}
	})
	t.Run("fail due to block", func(t *testing.T) {
		err := userStoreService.ValidatePassword(ctx, TestUser.ID, "invalid")
		if assert.Error(t, err) {
			if assert.EqualError(t, err, "password mismatch") {
				err = userStoreService.ValidatePassword(ctx, TestUser.ID, "invalid")
				if assert.Error(t, err) {
					assert.EqualError(t, err, "credential blocked")
				}
			}
		}
	})
	t.Run("reset password", func(t *testing.T) {
		err := userStoreService.SetPassword(ctx, TestUser.ID, "password")
		if assert.NoError(t, err) {
			err = userStoreService.ValidatePassword(ctx, TestUser.ID, "password")
			assert.NoError(t, err)
		}
	})
	userStoreService.RollbackTransaction(ctx)
}
