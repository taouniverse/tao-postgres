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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taouniverse/tao"
)

// TestIntegration 集成测试
// 通过环境变量 TAO_TEST_MULTI_INSTANCE 指定配置文件：
//   - 不设置或设置为空/"false"：使用默认单实例配置 (test.yaml)
//   - 设置为 "true"：使用多实例配置 (test_multi.yaml)
//
// 示例：
//
//	# 运行所有测试（单实例模式，默认）
//	go test -v ./...
//
//	# 运行所有测试（多实例模式）
//	TAO_TEST_MULTI_INSTANCE=true go test -v ./...
func TestIntegration(t *testing.T) {
	isMulti := os.Getenv("TAO_TEST_MULTI_INSTANCE") == "true"

	configPath := "./test.yaml"
	if isMulti {
		configPath = "./test_multi.yaml"
		t.Log("使用多实例配置进行测试")
	} else {
		t.Log("使用单实例配置进行测试")
	}

	err := tao.SetConfigPath(configPath)
	assert.Nil(t, err)

	// 测试获取默认实例
	db, err := DB()
	assert.Nil(t, err)

	sqlDB, err := db.DB()
	assert.Nil(t, err)

	err = sqlDB.Ping()
	assert.Nil(t, err)
	t.Log("默认实例连接成功")

	// 如果是多实例配置，测试获取指定实例
	if !isMulti {
		err = tao.Run(context.Background(), nil)
		assert.Nil(t, err)
		return
	}

	// 获取指定实例 master
	masterDB, err := GetDB("master")
	assert.Nil(t, err)
	assert.Equal(t, db, masterDB, "DB() 应该返回默认实例 master")
	t.Log("master 实例获取成功")

	sqlDB2, err := masterDB.DB()
	assert.Nil(t, err)

	err = sqlDB2.Ping()
	assert.Nil(t, err)
	t.Log("master 实例 Ping 成功")

	// 获取指定实例 replica
	replicaDB, err := GetDB("replica")
	assert.Nil(t, err)
	t.Log("replica 实例获取成功")

	sqlDB3, err := replicaDB.DB()
	assert.Nil(t, err)

	err = sqlDB3.Ping()
	assert.Nil(t, err)
	t.Log("replica 实例 Ping 成功")

	// 验证 master 和 replica 是不同的连接
	assert.NotEqual(t, masterDB, replicaDB, "master 和 replica 应该是不同的连接")
	t.Log("master 和 replica 是不同的连接")

	err = tao.Run(context.Background(), nil)
	assert.Nil(t, err)
}
