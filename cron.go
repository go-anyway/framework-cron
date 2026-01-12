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
	"sync"
	"time"

	"github.com/go-anyway/framework-log"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// Scheduler Cron 调度器
type Scheduler struct {
	cron    *cron.Cron
	tasks   map[cron.EntryID]*TaskInfo
	tasksMu sync.RWMutex
	opts    *Options
}

// NewScheduler 创建新的调度器
func NewScheduler(opts *Options) *Scheduler {
	if opts == nil {
		opts = DefaultOptions()
	}

	c := cron.New(
		cron.WithLocation(opts.Location),
		cron.WithSeconds(), // 支持秒级精度
	)

	return &Scheduler{
		cron:  c,
		tasks: make(map[cron.EntryID]*TaskInfo),
		opts:  opts,
	}
}

// AddTask 添加任务
func (s *Scheduler) AddTask(spec string, task Task) (cron.EntryID, error) {
	var taskID cron.EntryID
	id, err := s.cron.AddFunc(spec, func() {
		ctx := context.Background()
		startTime := time.Now()

		s.tasksMu.Lock()
		info := s.tasks[taskID]
		if info != nil {
			info.LastRun = startTime
			info.RunCount++
		}
		s.tasksMu.Unlock()

		// 如果启用了追踪，使用追踪包装
		var err error
		if s.opts.EnableTrace {
			err = runTaskWithTrace(ctx, task, func(ctx context.Context) error {
				return task.Run(ctx)
			})
		} else {
			err = task.Run(ctx)
		}
		duration := time.Since(startTime)

		s.tasksMu.Lock()
		if info != nil {
			info.NextRun = s.cron.Entry(taskID).Next
			if err != nil {
				info.ErrorCount++
				info.LastError = err
				log.FromContext(ctx).Error("Cron task failed",
					zap.String("task", task.Name()),
					zap.Duration("duration", duration),
					zap.Error(err),
				)
			} else {
				log.FromContext(ctx).Info("Cron task completed",
					zap.String("task", task.Name()),
					zap.Duration("duration", duration),
				)
			}
		}
		s.tasksMu.Unlock()
	})

	if err != nil {
		return 0, err
	}

	taskID = id

	entry := s.cron.Entry(id)
	s.tasksMu.Lock()
	s.tasks[id] = &TaskInfo{
		ID:      id,
		Name:    task.Name(),
		Spec:    spec,
		Task:    task,
		NextRun: entry.Next,
	}
	s.tasksMu.Unlock()

	return id, nil
}

// RemoveTask 移除任务
func (s *Scheduler) RemoveTask(id cron.EntryID) {
	s.cron.Remove(id)
	s.tasksMu.Lock()
	delete(s.tasks, id)
	s.tasksMu.Unlock()
}

// Start 启动调度器
func (s *Scheduler) Start() {
	s.cron.Start()
}

// Stop 停止调度器
func (s *Scheduler) Stop() context.Context {
	return s.cron.Stop()
}

// GetTaskInfo 获取任务信息
func (s *Scheduler) GetTaskInfo(id cron.EntryID) *TaskInfo {
	s.tasksMu.RLock()
	defer s.tasksMu.RUnlock()
	return s.tasks[id]
}

// GetAllTasks 获取所有任务信息
func (s *Scheduler) GetAllTasks() []*TaskInfo {
	s.tasksMu.RLock()
	defer s.tasksMu.RUnlock()

	tasks := make([]*TaskInfo, 0, len(s.tasks))
	for _, info := range s.tasks {
		// 更新 NextRun
		entry := s.cron.Entry(info.ID)
		info.NextRun = entry.Next
		tasks = append(tasks, info)
	}
	return tasks
}
