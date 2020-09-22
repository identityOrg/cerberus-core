package core

import (
	"context"
	"database/sql"
	"github.com/jinzhu/gorm"
)

const txKey = "gorm.tx"

func beginTransaction(ctx context.Context, db *gorm.DB) *gorm.DB {
	opt := &sql.TxOptions{
		ReadOnly: true,
	}
	return db.BeginTx(ctx, opt)
}

func rollbackTransaction(db *gorm.DB) {
	db.Rollback()
}
