package core

import (
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
	userStoreService := NewUserStoreService(TestDb, 3, 5*time.Minute)
	t.Run("de-activate", func(t *testing.T) {
		err := userStoreService.DeactivateUser(TestUser.ID)
		assert.NoError(t, err)
	})
	t.Run("activate", func(t *testing.T) {
		err := userStoreService.ActivateUser(TestUser.ID)
		assert.NoError(t, err)
	})
	t.Run("user not exists", func(t *testing.T) {
		err := userStoreService.ActivateUser(1000)
		if assert.Error(t, err) {
			assert.NotNil(t, err)
		}
	})
}

func TestUserStoreServiceImpl_FindAllUser(t *testing.T) {
	userStoreService := NewUserStoreService(TestDb, 3, 5*time.Minute)
	allUser, count, err := userStoreService.FindAllUser(0, 5)
	assert.Nil(t, err)
	if assert.Equal(t, uint(3), count) {
		assert.Equal(t, uint(1), allUser[0].ID)
	}
}

func TestUserStoreServiceImpl_ValidatePassword(t *testing.T) {
	userStoreService := NewUserStoreService(TestDb, 3, 5*time.Minute)
	t.Run("valid", func(t *testing.T) {
		err := userStoreService.ValidatePassword(1, "password")
		assert.Nil(t, err)
	})
	t.Run("invalid", func(t *testing.T) {
		err := userStoreService.ValidatePassword(1, "password1")
		assert.Error(t, err)
	})
	t.Run("user not exists", func(t *testing.T) {
		err := userStoreService.ValidatePassword(2000, "password")
		assert.Error(t, err)
	})
	t.Run("cred not found", func(t *testing.T) {
		err := userStoreService.ValidatePassword(TestNoCredUser2.ID, "password")
		if assert.Error(t, err) {
			assert.EqualError(t, err, "credential not found")
		}
	})
}

func TestUserStoreServiceImpl_FindUserByEmail(t *testing.T) {
	userStoreService := NewUserStoreService(TestDb, 3, 5*time.Minute)
	t.Run("found", func(t *testing.T) {
		foundUser, err := userStoreService.FindUserByEmail(TestUser.EmailAddress)
		assert.NoError(t, err)
		assert.Equal(t, TestUser.ID, foundUser.ID)
	})
	t.Run("not found", func(t *testing.T) {
		foundUser, err := userStoreService.FindUserByEmail("invalid@domain.com")
		assert.Error(t, err)
		assert.Nil(t, foundUser)
	})
}

func TestUserStoreServiceImpl_FindUserByUsername(t *testing.T) {
	userStoreService := NewUserStoreService(TestDb, 3, 5*time.Minute)
	t.Run("found", func(t *testing.T) {
		foundUser, err := userStoreService.FindUserByUsername(TestUser.Username)
		assert.NoError(t, err)
		assert.Equal(t, TestUser.ID, foundUser.ID)
	})
	t.Run("not found", func(t *testing.T) {
		foundUser, err := userStoreService.FindUserByUsername("invalid")
		assert.Error(t, err)
		assert.Nil(t, foundUser)
	})
}

func TestUserStoreServiceImpl_ValidateTOTP(t *testing.T) {
	userStoreService := NewUserStoreService(TestDb, 3, 5*time.Minute)
	t.Run("valid", func(t *testing.T) {
		code, err := totp.GenerateCode(TestUser.Credentials[1].Value, time.Now())
		if assert.NoError(t, err) {
			err = userStoreService.ValidateTOTP(TestUser.ID, code)
			assert.NoError(t, err)
		}
	})
	t.Run("invalid", func(t *testing.T) {
		err := userStoreService.ValidateTOTP(TestUser.ID, "code")
		assert.Error(t, err)
	})
	t.Run("no user", func(t *testing.T) {
		err := userStoreService.ValidateTOTP(2000, "code")
		assert.Error(t, err)
	})
}

func TestUserStoreServiceImpl_SetPassword(t *testing.T) {
	userStoreService := NewUserStoreService(TestDb, 3, 5*time.Minute)
	t.Run("success", func(t *testing.T) {
		err := userStoreService.SetPassword(TestNoCredUser.ID, "new password")
		assert.NoError(t, err)
	})
	t.Run("user not found", func(t *testing.T) {
		err := userStoreService.SetPassword(2000, "new password")
		assert.Error(t, err)
	})
}

func TestUserStoreServiceImpl_GenerateTOTP(t *testing.T) {
	userStoreService := NewUserStoreService(TestDb, 3, 5*time.Minute)
	t.Run("success", func(t *testing.T) {
		image, secret, err := userStoreService.GenerateTOTP(TestNoCredUser.ID, "cerberus")
		if assert.NoError(t, err) {
			assert.NotNil(t, image)
			assert.NotEqual(t, "", secret)
		}
	})
	t.Run("user not found", func(t *testing.T) {
		image, secret, err := userStoreService.GenerateTOTP(2000, "cerberus")
		if assert.Error(t, err) {
			assert.Nil(t, image)
			assert.EqualError(t, err, "user not found")
			assert.Equal(t, "", secret)
		}
	})
}

func TestUserStoreServiceImpl_Credential_Block_Unblock(t *testing.T) {
	userStoreService := NewUserStoreService(TestDb, 1, 5*time.Minute)
	t.Run("blocked", func(t *testing.T) {
		err := userStoreService.SetPassword(TestUser.ID, "other password")
		if assert.NoError(t, err) {
			err = userStoreService.ValidatePassword(TestUser.ID, "invalid")
			if assert.Error(t, err) {
				assert.EqualError(t, err, "password mismatch")
			}
		}
	})
	t.Run("fail due to block", func(t *testing.T) {
		err := userStoreService.ValidatePassword(TestUser.ID, "invalid")
		if assert.Error(t, err) {
			if assert.EqualError(t, err, "password mismatch") {
				err = userStoreService.ValidatePassword(TestUser.ID, "invalid")
				if assert.Error(t, err) {
					assert.EqualError(t, err, "credential blocked")
				}
			}
		}
	})
	t.Run("reset password", func(t *testing.T) {
		err := userStoreService.SetPassword(TestUser.ID, "password")
		if assert.NoError(t, err) {
			err = userStoreService.ValidatePassword(TestUser.ID, "password")
			assert.NoError(t, err)
		}
	})
}
