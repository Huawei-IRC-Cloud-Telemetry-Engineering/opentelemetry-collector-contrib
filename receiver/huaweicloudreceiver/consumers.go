// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package huaweicloudreceiver

import (
	"context"

	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
)

type receiverProcessor interface {
	processReceivedData(ctx context.Context, data []byte) error
}

type metricsReceiver struct {
	consumer consumer.Metrics
}

func (r *metricsReceiver) processReceivedData(ctx context.Context, data []byte) error {
	unmarshaler := &pmetric.JSONUnmarshaler{}

	metrics, err := unmarshaler.UnmarshalMetrics(data)
	if err != nil {
		return err
	}
	return r.consumer.ConsumeMetrics(ctx, metrics)
}

func newHuaweiCloudMetricsReceiver(cfg *Config, metrics consumer.Metrics, settings receiver.Settings) *huaweiReceiver {
	return newHuaweiCloudReceiver(settings, Metrics, cfg, &metricsReceiver{consumer: metrics})
}

type logsReceiver struct {
	consumer consumer.Logs
}

func (r *logsReceiver) processReceivedData(ctx context.Context, data []byte) error {
	unmarshaler := &plog.JSONUnmarshaler{}

	logs, err := unmarshaler.UnmarshalLogs(data)
	if err != nil {
		return err
	}
	return r.consumer.ConsumeLogs(ctx, logs)
}

func newHuaweiCloudLogsReceiver(cfg *Config, logs consumer.Logs, settings receiver.Settings) *huaweiReceiver {
	return newHuaweiCloudReceiver(settings, Logs, cfg, &logsReceiver{consumer: logs})
}
