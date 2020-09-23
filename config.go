package core

import "time"

type Config struct {
	//DatabaseDialect        string
	//DataSourceName         string
	EncryptionKey          string
	MaxInvalidLoginAttempt uint
	InvalidAttemptWindow   time.Duration
	TOTPSecretLength       uint
}
