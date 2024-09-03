// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package huaweicloudreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/huaweicloudreceiver"

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/receiver/receivertest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/golden"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest/pmetrictest"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/huaweicloudreceiver/internal/mocks"
)

func TestHuaweiCloudMetricsReceiverIntegration(t *testing.T) {
	mc := mocks.NewCesClient(t)

	mc.On("ListMetrics", mock.Anything).Return(&model.ListMetricsResponse{
		Metrics: &[]model.MetricInfoList{
			{
				Namespace:  "SYS.ECS",
				MetricName: "cpu_util",
				Dimensions: []model.MetricsDimension{
					{
						Name:  "instance_id",
						Value: "faea5b75-e390-4e2b-8733-9226a9026070",
					},
				},
				Unit: "%",
			},
		},
	}, nil)

	mc.On("ShowMetricData", mock.Anything).Return(&model.ShowMetricDataResponse{
		MetricName: stringPtr("cpu_util"),
		Datapoints: &[]model.Datapoint{
			{
				Average:   float64Ptr(45.67),
				Timestamp: 1556625610000,
			},
			{
				Average:   float64Ptr(89.01),
				Timestamp: 1556625715000,
			},
		},
	}, nil)

	sink := &consumertest.MetricsSink{}
	cfg := createDefaultConfig().(*Config)
	cfg.RegionName = "us-east-2"
	cfg.CollectionInterval = time.Second
	cfg.ProjectID = "my-project"
	cfg.Filter = "average"

	recv, err := NewFactory().CreateMetricsReceiver(
		context.Background(),
		receivertest.NewNopSettings(),
		cfg,
		sink,
	)
	require.NoError(t, err)

	rcvr, ok := recv.(*huaweiReceiver)
	require.True(t, ok)
	rcvr.clients = &clients{ces: mc}

	err = recv.Start(context.Background(), componenttest.NewNopHost())
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		return sink.DataPointCount() > 0
	}, 5*time.Second, 10*time.Millisecond)

	err = recv.Shutdown(context.Background())
	require.NoError(t, err)

	metrics := sink.AllMetrics()[0]

	expectedMetrics, err := golden.ReadMetrics(filepath.Join("testdata", "golden", "metrics_golden.yaml"))
	require.NoError(t, err)
	require.NoError(t, pmetrictest.CompareMetrics(expectedMetrics, metrics))
}
