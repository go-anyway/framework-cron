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
	"context"
	"time"

	"github.com/robfig/cron/v3"
)

// Task 定时任务接口
type Task interface {
	// Name 返回任务名称
	Name() string
	// Run 执行任务
	Run(ctx context.Context) error
}

// TaskFunc 任务函数类型
type TaskFunc func(ctx context.Context) error

// Name 返回任务名称
func (f TaskFunc) Name() string {
	return "anonymous"
}

// Run 执行任务
func (f TaskFunc) Run(ctx context.Context) error {
	return f(ctx)
}

// NamedTask 命名任务
type NamedTask struct {
	name string
	fn   TaskFunc
}

// NewNamedTask 创建命名任务
func NewNamedTask(name string, fn TaskFunc) *NamedTask {
	return &NamedTask{
		name: name,
		fn:   fn,
	}
}

// Name 返回任务名称
func (t *NamedTask) Name() string {
	return t.name
}

// Run 执行任务
func (t *NamedTask) Run(ctx context.Context) error {
	return t.fn(ctx)
}

// TaskInfo 任务信息
type TaskInfo struct {
	ID         cron.EntryID
	Name       string
	Spec       string
	Task       Task
	LastRun    time.Time
	NextRun    time.Time
	RunCount   int64
	ErrorCount int64
	LastError  error
}
