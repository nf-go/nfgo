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

package nmetrics

import (
	"github.com/nf-go/nfgo/nconf"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

func (s *server) regitserDBCollector(config *nconf.Config) error {
	if config.DB != nil && s.opts != nil && s.opts.db != nil {
		sqldb, err := s.opts.db.DB()
		if err != nil {
			return err
		}
		return s.registry.Register(collectors.NewDBStatsCollector(sqldb, config.DB.Database))
	}
	return nil
}
