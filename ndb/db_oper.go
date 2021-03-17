// Copyright 2021 The nfgo Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ndb

import (
	"context"
	"database/sql"
	"fmt"

	"gorm.io/gorm"
)

type ctxKeyDb struct{}

// DBOper -
type DBOper interface {
	DB(ctx context.Context) *gorm.DB
	Transactional(ctx context.Context, fn func(ctx context.Context) error) (err error)
}

// NewDBOper -
func NewDBOper(db *gorm.DB) DBOper {
	return &dbOperImpl{
		db: db,
	}
}

type dbOperImpl struct {
	db *gorm.DB
}

func (o *dbOperImpl) DB(ctx context.Context) *gorm.DB {
	v := ctx.Value(ctxKeyDb{})
	if dbInCtx, ok := v.(*gorm.DB); ok {
		return dbInCtx.WithContext(ctx)
	}
	return o.db.WithContext(ctx)
}

func (o *dbOperImpl) Transactional(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	panicked := true
	tx := o.db.Begin(&sql.TxOptions{Isolation: sql.LevelDefault})
	if tx.Error != nil {
		return fmt.Errorf("unable to begin transaction: %w", tx.Error)
	}

	defer func() {
		// Make sure to rollback when panic, Block error or Commit error
		if panicked || err != nil {
			tx.Rollback()
		}
	}()

	ctx = context.WithValue(ctx, ctxKeyDb{}, tx)
	err = fn(ctx)

	if err == nil {
		err = tx.Commit().Error
	}
	panicked = false
	return
}
