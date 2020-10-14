package core

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
)

const txKey = "gorm.tx"

func beginTransaction(ctx context.Context, db *gorm.DB) *gorm.DB {
	opt := &sql.TxOptions{
		ReadOnly: true,
	}
	return db.Begin(opt)
}

func rollbackTransaction(db *gorm.DB) {
	db.Rollback()
}
