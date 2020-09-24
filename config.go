package core

import "time"

type Config struct {
	EncryptionKey          string
	MaxInvalidLoginAttempt uint
	InvalidAttemptWindow   time.Duration
	TOTPSecretLength       uint
	PasswordCost           int
}
