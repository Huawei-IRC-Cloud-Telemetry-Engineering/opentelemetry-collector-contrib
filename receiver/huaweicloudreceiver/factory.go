// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package huaweicloudreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/huaweicloudreceiver"

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/huaweicloudreceiver/internal/metadata"
)

func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		metadata.Type,
		createDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, metadata.MetricsStability),
		receiver.WithLogs(createLogsReceiver, metadata.LogsStability),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		BackOffConfig: configretry.BackOffConfig{
			Enabled:             true,
			InitialInterval:     100 * time.Millisecond,
			MaxInterval:         time.Second,
			MaxElapsedTime:      15 * time.Second,
			RandomizationFactor: backoff.DefaultRandomizationFactor,
			Multiplier:          backoff.DefaultMultiplier,
		},
		HuaweiSessionConfig: HuaweiSessionConfig{
			NoVerifySSL: false,
		},
	}
}

func createMetricsReceiver(
	_ context.Context,
	params receiver.Settings,
	cfg component.Config,
	next consumer.Metrics) (receiver.Metrics, error) {
	return newHuaweiCloudMetricsReceiver(cfg.(*Config), next, params), nil

}

func createLogsReceiver(
	_ context.Context,
	params receiver.Settings,
	cfg component.Config,
	next consumer.Logs) (receiver.Logs, error) {
	return newHuaweiCloudLogsReceiver(cfg.(*Config), next, params), nil
}
