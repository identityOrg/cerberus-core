package core

import (
	"context"
	"database/sql"
	"github.com/jinzhu/gorm"
)

const txKey = "gorm.tx"

func beginTransaction(ctx context.Context, readOnly bool, db *gorm.DB) context.Context {
	opt := &sql.TxOptions{
		ReadOnly: readOnly,
	}
	tx := db.BeginTx(ctx, opt)
	return context.WithValue(ctx, txKey, tx)
}

func commitTransaction(ctx context.Context) context.Context {
	tx := ctx.Value(txKey)
	if tx != nil {
		if db, ok := tx.(*gorm.DB); ok {
			db.Commit()
		}
	}
	return context.WithValue(ctx, txKey, nil)
}

func rollbackTransaction(ctx context.Context) context.Context {
	tx := ctx.Value(txKey)
	if tx != nil {
		if db, ok := tx.(*gorm.DB); ok {
			db.Rollback()
		}
	}
	return context.WithValue(ctx, txKey, nil)
}

func getTransaction(ctx context.Context) *gorm.DB {
	tx := ctx.Value(txKey)
	if tx != nil {
		if db, ok := tx.(*gorm.DB); ok {
			return db
		}
	}
	panic("transaction doesnt exists")
}
