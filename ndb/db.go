package ndb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"nfgo.ga/nfgo/nconf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"nfgo.ga/nfgo/ncontext"
)

// NewDB -
func NewDB(dbConfig *nconf.DbConfig) (*gorm.DB, error) {
	if dbConfig == nil {
		return nil, errors.New("dbConfig is nil")
	}

	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger:                 newLogger(dbConfig),
		PrepareStmt:            dbConfig.PrepareStmt,
		SkipDefaultTransaction: dbConfig.SkipDefaultTransaction,
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
		dbConfig.Username, dbConfig.Password, dbConfig.Host,
		dbConfig.Port, dbConfig.Database, dbConfig.Charset)

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("fail to open db: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(int(dbConfig.MaxIdle))
	sqlDB.SetMaxOpenConns(int(dbConfig.MaxOpen))

	return db, nil
}

// WithContext -
func WithContext(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	v := ctx.Value(ncontext.CtxKeyDb)
	if dbInCtx, ok := v.(*gorm.DB); ok {
		return dbInCtx.WithContext(ctx)
	}
	return defaultDB.WithContext(ctx)
}

// Transactional -
func Transactional(ctx context.Context, db *gorm.DB, fn func(ctx context.Context) error) (err error) {
	panicked := true
	tx := db.Begin(&sql.TxOptions{Isolation: sql.LevelDefault})
	if tx.Error != nil {
		return fmt.Errorf("Unable to begin transaction: %w", tx.Error)
	}

	defer func() {
		// Make sure to rollback when panic, Block error or Commit error
		if panicked || err != nil {
			tx.Rollback()
		}
	}()

	ctx = context.WithValue(ctx, ncontext.CtxKeyDb, tx)
	err = fn(ctx)

	if err == nil {
		err = tx.Commit().Error
	}
	panicked = false
	return
}
