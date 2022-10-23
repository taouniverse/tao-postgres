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
	"fmt"
	"github.com/taouniverse/tao"
	// Load the required dependencies.
	// An error occurs when there was no package in the root directory.
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/**
import _ "github.com/taouniverse/tao-postgres"
*/

// P config of postgres
var P = new(Config)

func init() {
	err := tao.Register(ConfigKey, P, setup)
	if err != nil {
		panic(err.Error())
	}
}

// DB orm client of mysql
var DB *gorm.DB

// setup unit with the global config 'P'
// execute when init tao universe
func setup() (err error) {
	// FIXME dbname must be last one
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s sslmode=%s TimeZone=%s dbname=%s",
		P.Host, P.Port, P.User, P.Password, P.SSL, P.TimeZone, P.DB,
	)

	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		return tao.NewErrorWrapped("postgres: fail to create gorm client", err)
	}

	return nil
}
