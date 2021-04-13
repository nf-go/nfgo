// Copyright 2020 The nfgo Authors. All Rights Reserved.
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
	"errors"
	"fmt"

	"nfgo.ga/nfgo/nconf"
	"nfgo.ga/nfgo/nlog"
	"nfgo.ga/nfgo/nutil/ntypes"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type dbOptions struct {
	dialector gorm.Dialector
}

// DBOption -
type DBOption func(*dbOptions)

// DialectorOption -
func DialectorOption(dialector gorm.Dialector) DBOption {
	return func(opts *dbOptions) {
		opts.dialector = dialector
	}
}

// NewDB -
func NewDB(dbConfig *nconf.DbConfig, opt ...DBOption) (*gorm.DB, error) {
	if dbConfig == nil {
		return nil, errors.New("dbConfig is nil")
	}

	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger:                 newLogger(dbConfig),
		PrepareStmt:            ntypes.BoolValue(dbConfig.PrepareStmt),
		SkipDefaultTransaction: ntypes.BoolValue(dbConfig.SkipDefaultTransaction),
	}

	opts := &dbOptions{}
	for _, o := range opt {
		o(opts)
	}
	if ntypes.IsNil(opts.dialector) {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
			dbConfig.Username, dbConfig.Password, dbConfig.Host,
			dbConfig.Port, dbConfig.Database, dbConfig.Charset)
		opts.dialector = mysql.Open(dsn)
	}

	db, err := gorm.Open(opts.dialector, gormConfig)
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

// MustNewDB -
func MustNewDB(dbConfig *nconf.DbConfig, opt ...DBOption) *gorm.DB {
	db, err := NewDB(dbConfig, opt...)
	if err != nil {
		nlog.Fatal("fail to new db: ", err)
	}
	return db
}
