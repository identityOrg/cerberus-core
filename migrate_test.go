package core

import (
	"github.com/identityOrg/oidcsdk"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestMigrateDB(t *testing.T) {
	_ = os.Remove("migrate.db")
	db, err := gorm.Open("sqlite3", "migrate.db")
	if assert.NoError(t, err) {
		err := SetupDBStructure(db, true, true)
		if assert.NoError(t, err) {
			config := &Config{
				EncryptionKey:          "ewewevwev",
				MaxInvalidLoginAttempt: 3,
				InvalidAttemptWindow:   5 * time.Minute,
				TOTPSecretLength:       6,
				PasswordCost:           8,
			}
			sdkConfig := oidcsdk.NewConfig("http://localhost:8080")
			err := SetupDemoData(db, config, sdkConfig, "")
			assert.NoError(t, err)
		}
	}
}
