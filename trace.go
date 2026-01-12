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

	"github.com/go-anyway/framework-log"
	pkgtrace "github.com/go-anyway/framework-trace"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// runTaskWithTrace 带追踪的任务执行包装器
func runTaskWithTrace(
	ctx context.Context,
	task Task,
	handler func(context.Context) error,
) error {
	startTime := time.Now()

	// 创建追踪 span
	ctx, span := pkgtrace.StartSpan(ctx, "cron.task",
		trace.WithAttributes(
			attribute.String("cron.task.name", task.Name()),
		),
	)
	defer span.End()

	// 执行任务
	err := handler(ctx)
	duration := time.Since(startTime)

	// 记录日志
	if err != nil {
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

	// 更新追踪状态
	span.SetAttributes(
		attribute.Float64("cron.duration_ms", float64(duration.Milliseconds())),
	)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return err
}
