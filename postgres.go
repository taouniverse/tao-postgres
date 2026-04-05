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
	"fmt"

	"github.com/taouniverse/tao"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/**
import _ "github.com/taouniverse/tao-postgres"
*/

// P is the global config instance for tao-postgres
var P = &Config{}

// Factory is the global factory instance for managing gorm.DB
var Factory *tao.BaseFactory[*gorm.DB]

func init() {
	var err error
	Factory, err = tao.Register(ConfigKey, P, NewPostgres)
	if err != nil {
		panic(err.Error())
	}
}

// NewPostgres creates a new PostgreSQL client for factory pattern
func NewPostgres(name string, config InstanceConfig) (*gorm.DB, func() error, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s sslmode=%s TimeZone=%s dbname=%s",
		config.Host, config.Port, config.User, config.Password, config.SSL, config.TimeZone, config.DB,
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, tao.NewErrorWrapped("postgres: fail to create gorm client", err)
	}

	closer := func() error {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}

	return db, closer, nil
}

// DB returns the default gorm DB instance
func DB() (*gorm.DB, error) {
	return Factory.Get(P.GetDefaultInstanceName())
}

// GetDB returns the gorm DB instance by name
func GetDB(name string) (*gorm.DB, error) {
	return Factory.Get(name)
}
