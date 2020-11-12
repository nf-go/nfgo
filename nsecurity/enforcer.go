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

package nsecurity

import (
	"fmt"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
	"nfgo.ga/nfgo/nconf"
	"nfgo.ga/nfgo/nlog"
)

// NewEnforcer -
func NewEnforcer(securityConfig *nconf.SecurityConfig, db *gorm.DB) (casbin.IEnforcer, error) {
	adpt, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("fail to create casbin gorm-adapter: %w", err)
	}

	m := model.NewModel()
	if err := m.LoadModelFromText(securityConfig.Model); err != nil {
		return nil, fmt.Errorf("fail to load model from seucrity config: %w", err)
	}

	enforcer, err := casbin.NewEnforcer(m, adpt)
	if err != nil {
		return nil, fmt.Errorf("fail to create enforcer: %w", err)
	}
	return enforcer, nil
}

// MustNewEnforcer -
func MustNewEnforcer(securityConfig *nconf.SecurityConfig, db *gorm.DB) casbin.IEnforcer {
	enforcer, err := NewEnforcer(securityConfig, db)
	if err != nil {
		nlog.Fatal("fail to create enforcer: ", err)
	}
	return enforcer
}

// InitPolicy -
func InitPolicy(enforcer casbin.IEnforcer, securityConfig *nconf.SecurityConfig, rules [][]string) error {
	// load from db
	if err := enforcer.LoadPolicy(); err != nil {
		return err
	}
	// clear all
	enforcer.ClearPolicy()

	// add policies from config
	for _, anno := range securityConfig.Anons {
		if _, err := enforcer.AddNamedPolicy("p", "anonymous", anno, "*"); err != nil {
			return err
		}
	}
	for _, policy := range securityConfig.Policies {
		ps := strings.Split(policy, ",")
		if len(ps) > 2 {
			params := make([]interface{}, len(ps))
			for i := range ps {
				params[i] = strings.TrimSpace(ps[i])
			}
			enforcer.AddNamedPolicy(ps[0], params[1:]...)
		}
	}

	// save policies to db
	if err := enforcer.SavePolicy(); err != nil {
		return err
	}

	if _, err := enforcer.AddPolicies(rules); err != nil {
		return err
	}

	return nil
}
