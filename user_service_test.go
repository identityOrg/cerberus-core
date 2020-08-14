package core

import (
	"encoding/base32"
	"github.com/pquerna/otp/totp"
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
