package core

import (
	"context"
	"encoding/base32"
	"fmt"
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
	config := &Config{
		MaxInvalidLoginAttempt: 3,
		InvalidAttemptWindow:   5 * time.Minute,
		TOTPSecretLength:       6,
	}
	userStoreService := NewUserStoreServiceImpl(TestDb, config)
	userStoreService.Db = beginTransaction(ctx, userStoreService.Db)
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
	rollbackTransaction(userStoreService.Db)
}

func TestUserStoreServiceImpl_FindAllUser(t *testing.T) {
	ctx := context.Background()
	config := &Config{
		MaxInvalidLoginAttempt: 3,
		InvalidAttemptWindow:   5 * time.Minute,
		TOTPSecretLength:       6,
	}
	userStoreService := NewUserStoreServiceImpl(TestDb, config)
	userStoreService.Db = beginTransaction(ctx, userStoreService.Db)
	allUser, count, err := userStoreService.FindAllUser(ctx, 0, 5)
	assert.Nil(t, err)
	if assert.Equal(t, uint(3), count) {
		assert.Equal(t, uint(1), allUser[0].ID)
	}
	rollbackTransaction(userStoreService.Db)
}

func TestUserStoreServiceImpl_ValidatePassword(t *testing.T) {
	ctx := context.Background()
	config := &Config{
		MaxInvalidLoginAttempt: 3,
		InvalidAttemptWindow:   5 * time.Minute,
		TOTPSecretLength:       6,
	}
	userStoreService := NewUserStoreServiceImpl(TestDb, config)
	userStoreService.Db = beginTransaction(ctx, userStoreService.Db)
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
			assert.EqualError(t, err, fmt.Sprintf("password not set for user %d", TestNoCredUser2.ID))
		}
	})
	rollbackTransaction(userStoreService.Db)
}

func TestUserStoreServiceImpl_FindUserByEmail(t *testing.T) {
	ctx := context.Background()
	config := &Config{
		MaxInvalidLoginAttempt: 3,
		InvalidAttemptWindow:   5 * time.Minute,
		TOTPSecretLength:       6,
	}
	userStoreService := NewUserStoreServiceImpl(TestDb, config)
	userStoreService.Db = beginTransaction(ctx, userStoreService.Db)
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
	rollbackTransaction(userStoreService.Db)
}

func TestUserStoreServiceImpl_FindUserByUsername(t *testing.T) {
	ctx := context.Background()
	config := &Config{
		MaxInvalidLoginAttempt: 3,
		InvalidAttemptWindow:   5 * time.Minute,
		TOTPSecretLength:       6,
	}
	userStoreService := NewUserStoreServiceImpl(TestDb, config)
	userStoreService.Db = beginTransaction(ctx, userStoreService.Db)
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
	rollbackTransaction(userStoreService.Db)
}

func TestUserStoreServiceImpl_ValidateTOTP(t *testing.T) {
	ctx := context.Background()
	config := &Config{
		MaxInvalidLoginAttempt: 3,
		InvalidAttemptWindow:   5 * time.Minute,
		TOTPSecretLength:       6,
	}
	userStoreService := NewUserStoreServiceImpl(TestDb, config)
	userStoreService.Db = beginTransaction(ctx, userStoreService.Db)
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
	rollbackTransaction(userStoreService.Db)
}

func TestUserStoreServiceImpl_SetPassword(t *testing.T) {
	ctx := context.Background()
	config := &Config{
		MaxInvalidLoginAttempt: 3,
		InvalidAttemptWindow:   5 * time.Minute,
		TOTPSecretLength:       6,
	}
	userStoreService := NewUserStoreServiceImpl(TestDb, config)
	userStoreService.Db = beginTransaction(ctx, userStoreService.Db)
	t.Run("success", func(t *testing.T) {
		err := userStoreService.SetPassword(ctx, TestNoCredUser.ID, "new password")
		assert.NoError(t, err)
	})
	t.Run("user not found", func(t *testing.T) {
		err := userStoreService.SetPassword(ctx, 2000, "new password")
		assert.Error(t, err)
	})
	rollbackTransaction(userStoreService.Db)
}

func TestUserStoreServiceImpl_GenerateTOTP(t *testing.T) {
	ctx := context.Background()
	config := &Config{
		MaxInvalidLoginAttempt: 3,
		InvalidAttemptWindow:   5 * time.Minute,
		TOTPSecretLength:       6,
	}
	userStoreService := NewUserStoreServiceImpl(TestDb, config)
	userStoreService.Db = beginTransaction(ctx, userStoreService.Db)
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
			assert.EqualError(t, err, "user not found with id 2000")
			assert.Equal(t, "", secret)
		}
	})
	rollbackTransaction(userStoreService.Db)
}

func TestUserStoreServiceImpl_Credential_Block_Unblock(t *testing.T) {
	ctx := context.Background()
	config := &Config{
		MaxInvalidLoginAttempt: 1,
		InvalidAttemptWindow:   5 * time.Minute,
		TOTPSecretLength:       6,
	}
	userStoreService := NewUserStoreServiceImpl(TestDb, config)
	userStoreService.Db = beginTransaction(ctx, userStoreService.Db)
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
	rollbackTransaction(userStoreService.Db)
}

func TestUserStoreServiceImpl_GetUser(t *testing.T) {
	ctx := context.Background()
	config := &Config{
		MaxInvalidLoginAttempt: 1,
		InvalidAttemptWindow:   5 * time.Minute,
		TOTPSecretLength:       6,
	}
	userStoreService := NewUserStoreServiceImpl(TestDb, config)
	userStoreService.Db = beginTransaction(ctx, userStoreService.Db)
	t.Run("success", func(t *testing.T) {
		user, err := userStoreService.GetUser(ctx, TestUser.ID)
		if assert.NoError(t, err) {
			if assert.NotNil(t, user) {
				assert.Equal(t, TestUser.Username, user.Username)
			}
		}
	})
	t.Run("fail", func(t *testing.T) {
		user, err := userStoreService.GetUser(ctx, 2000)
		if assert.Error(t, err) {
			assert.Nil(t, user)
			assert.EqualError(t, err, "user not found with id 2000")
		}
	})
	rollbackTransaction(userStoreService.Db)
}

func TestUserStoreServiceImpl_GetClaims(t *testing.T) {
	ctx := context.Background()
	config := &Config{
		MaxInvalidLoginAttempt: 1,
		InvalidAttemptWindow:   5 * time.Minute,
		TOTPSecretLength:       6,
	}
	userStoreService := NewUserStoreServiceImpl(TestDb, config)
	userStoreService.Db = beginTransaction(ctx, userStoreService.Db)
	_, _ = userStoreService.GetClaims(ctx, "us", []string{"openid"}, []string{})
	rollbackTransaction(userStoreService.Db)
}
