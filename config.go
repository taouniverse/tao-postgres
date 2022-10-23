// Copyright 2022 huija
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package postgres

import (
	"context"
	"github.com/taouniverse/tao"
)

// ConfigKey for this repo
const ConfigKey = "postgres"

// Config implements tao.Config
type Config struct {
	Host      string   `json:"host"`
	Port      int      `json:"port"`
	User      string   `json:"user"`
	Password  string   `json:"password"`
	DB        string   `json:"db"`
	SSL       string   `json:"ssl"`
	TimeZone  string   `json:"time_zone"`
	RunAfters []string `json:"run_after,omitempty"`
}

var defaultPostgres = &Config{
	Host:      "localhost",
	Port:      5432,
	User:      "tao",
	Password:  "123456qwe",
	SSL:       "disable",
	TimeZone:  "Asia/Shanghai",
	RunAfters: []string{},
}

// Name of Config
func (p *Config) Name() string {
	return ConfigKey
}

// ValidSelf with some default values
func (p *Config) ValidSelf() {
	if p.Host == "" {
		p.Host = defaultPostgres.Host
	}
	if p.Port == 0 {
		p.Port = defaultPostgres.Port
	}
	if p.User == "" {
		p.User = defaultPostgres.User
	}
	if p.Password == "" {
		p.Password = defaultPostgres.Password
	}
	if p.SSL == "" {
		p.SSL = defaultPostgres.SSL
	}
	if p.TimeZone == "" {
		p.TimeZone = defaultPostgres.TimeZone
	}
	if p.RunAfters == nil {
		p.RunAfters = defaultPostgres.RunAfters
	}
}

// ToTask transform itself to Task
func (p *Config) ToTask() tao.Task {
	return tao.NewTask(
		ConfigKey,
		func(ctx context.Context, param tao.Parameter) (tao.Parameter, error) {
			// non-block check
			select {
			case <-ctx.Done():
				return param, tao.NewError(tao.ContextCanceled, "%s: context has been canceled", ConfigKey)
			default:
			}
			// JOB code run after RunAfters, you can just do nothing here
			db, err := DB.DB()
			if err != nil {
				return param, err
			}

			err = db.Ping()
			return param, err
		})
}

// RunAfter defines pre task names
func (p *Config) RunAfter() []string {
	return p.RunAfters
}
