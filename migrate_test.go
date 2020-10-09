package core

import (
	"github.com/identityOrg/oidcsdk"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMigrateDB(t *testing.T) {
	_ = os.Remove("migrate.db")
	db, err := gorm.Open("sqlite3", "migrate.db")
	config := oidcsdk.NewConfig("http://localhost:8080")
	if assert.NoError(t, err) {
		defer db.Close()
		err := MigrateDB(db, &Config{}, config, false, false)
		assert.NoError(t, err)
	}
}
