// Copyright 2025 zampo.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// @contact  zampo3380@gmail.com

package cron

import (
	"fmt"
	"time"
)

// Config Cron 配置结构体（用于从配置文件创建）
type Config struct {
	Enabled     bool   `yaml:"enabled" env:"CRON_ENABLED" default:"true"`
	Location    string `yaml:"location" env:"CRON_LOCATION" default:"Local"`
	EnableTrace bool   `yaml:"enable_trace" env:"CRON_ENABLE_TRACE" default:"true"`
}

// Validate 验证 Cron 配置
func (c *Config) Validate() error {
	if c == nil {
		return fmt.Errorf("cron config cannot be nil")
	}
	if !c.Enabled {
		return nil // 如果未启用，不需要验证
	}
	return nil
}

// ToOptions 转换为 Options
func (c *Config) ToOptions() (*Options, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	if !c.Enabled {
		return nil, fmt.Errorf("cron is not enabled")
	}

	location := time.Local
	if c.Location != "" && c.Location != "Local" {
		var err error
		location, err = time.LoadLocation(c.Location)
		if err != nil {
			return nil, fmt.Errorf("invalid location: %w", err)
		}
	}

	return &Options{
		Location:    location,
		EnableTrace: c.EnableTrace,
	}, nil
}

// Options Cron 配置选项（内部使用）
type Options struct {
	Location    *time.Location // 时区
	EnableTrace bool           // 是否启用追踪
}

// DefaultOptions 返回默认配置选项
func DefaultOptions() *Options {
	return &Options{
		Location:    time.Local,
		EnableTrace: true,
	}
}
