// Copyright 2021-2026 huija
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

// InstanceConfig 单实例配置
type InstanceConfig struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	DB       string `json:"db" yaml:"db"`
	SSL      string `json:"ssl" yaml:"ssl"`
	TimeZone string `json:"time_zone" yaml:"time_zone"`
}

// Config 总配置，实现 tao.MultiConfig 接口
type Config struct {
	tao.BaseMultiConfig[InstanceConfig]
	RunAfters []string `json:"run_after,omitempty" yaml:"run_after,omitempty"`
}

var defaultInstance = &InstanceConfig{
	Host:     "localhost",
	Port:     5432,
	User:     "tao",
	Password: "123456qwe",
	SSL:      "disable",
	TimeZone: "Asia/Shanghai",
}

// Name of Config
func (c *Config) Name() string {
	return ConfigKey
}

// ValidSelf with some default values
func (c *Config) ValidSelf() {
	for i := range c.Instances {
		instance := &c.Instances[i].Cfg
		if instance.Host == "" {
			instance.Host = defaultInstance.Host
		}
		if instance.Port == 0 {
			instance.Port = defaultInstance.Port
		}
		if instance.User == "" {
			instance.User = defaultInstance.User
		}
		if instance.Password == "" {
			instance.Password = defaultInstance.Password
		}
		if instance.SSL == "" {
			instance.SSL = defaultInstance.SSL
		}
		if instance.TimeZone == "" {
			instance.TimeZone = defaultInstance.TimeZone
		}
	}
}

// ToTask transform itself to Task
func (c *Config) ToTask() tao.Task {
	return tao.NewTask(
		ConfigKey,
		func(ctx context.Context, param tao.Parameter) (tao.Parameter, error) {
			select {
			case <-ctx.Done():
				return param, tao.NewError(tao.ContextCanceled, "%s: context has been canceled", ConfigKey)
			default:
			}
			for _, inst := range c.Instances {
				name := inst.Name
				db, err := Factory.Get(name)
				if err != nil {
					return param, err
				}
				sqlDB, err := db.DB()
				if err != nil {
					return param, err
				}
				if err := sqlDB.Ping(); err != nil {
					return param, err
				}
			}
			return param, nil
		})
}

// RunAfter defines pre task names
func (c *Config) RunAfter() []string {
	return c.RunAfters
}
