// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"testing"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func TestConvertCESMetricsToOTLP(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]MetricData
		expected []byte
	}{
		{
			name: "Valid Metric Conversion",
			input: map[string]MetricData{
				"cpu_util": {
					MetricName: "cpu_util",
					Namespace:  "SYS.ECS",
					Dimensions: []model.MetricsDimension{
						{
							Name:  "instance_id",
							Value: "faea5b75-e390-4e2b-8733-9226a9026070",
						},
					},
					Datapoints: []model.Datapoint{
						{
							Average:   float64Ptr(0.5),
							Timestamp: 1556625610000,
						},
						{
							Average:   float64Ptr(0.7),
							Timestamp: 1556625715000,
						},
					},
					Unit: "%",
				},
				"network_vm_connections": {
					MetricName: "network_vm_connections",
					Namespace:  "SYS.ECS",
					Dimensions: []model.MetricsDimension{
						{
							Name:  "instance_id",
							Value: "06b4020f-461a-4a52-84da-53fa71c2f42b",
						},
					},
					Datapoints: []model.Datapoint{
						{
							Average:   float64Ptr(1),
							Timestamp: 1556625612000,
						},
						{
							Average:   float64Ptr(3),
							Timestamp: 1556625717000,
						},
					},
					Unit: "count",
				},
			},
			expected: expectedMetrics(t),
		},
		{
			name:     "Empty Metric Data",
			input:    map[string]MetricData{},
			expected: expectedEmptyMetrics(t),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := ConvertCESMetricsToOTLP("project_1", "eu-west-101", "average", tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, data)
		})
	}
}

func float64Ptr(f float64) *float64 {
	return &f
}

func expectedMetrics(t *testing.T) []byte {
	metrics := pmetric.NewMetrics()
	resourceMetric := metrics.ResourceMetrics().AppendEmpty()

	resource := resourceMetric.Resource()
	resource.Attributes().PutStr("cloud.provider", "huawei_cloud")
	resource.Attributes().PutStr("project.id", "project_1")
	resource.Attributes().PutStr("region", "eu-west-101")

	scopedMetric := resourceMetric.ScopeMetrics().AppendEmpty()
	scopedMetric.Scope().SetName("huawei_cloud_ces")
	scopedMetric.Scope().SetVersion("v1")

	metric := scopedMetric.Metrics().AppendEmpty()
	metric.SetName("cpu_util")
	metric.SetUnit("%")
	metric.Metadata().PutStr("instance_id", "faea5b75-e390-4e2b-8733-9226a9026070")
	metric.Metadata().PutStr("service.namespace", "SYS.ECS")

	dataPoints := metric.SetEmptyGauge().DataPoints()
	dp := dataPoints.AppendEmpty()
	dp.SetTimestamp(pcommon.NewTimestampFromTime(time.UnixMilli(1556625610000)))
	dp.SetDoubleValue(0.5)

	dp = dataPoints.AppendEmpty()
	dp.SetTimestamp(pcommon.NewTimestampFromTime(time.UnixMilli(1556625715000)))
	dp.SetDoubleValue(0.7)

	scopedMetric = resourceMetric.ScopeMetrics().AppendEmpty()
	scopedMetric.Scope().SetName("huawei_cloud_ces")
	scopedMetric.Scope().SetVersion("v1")

	metric = scopedMetric.Metrics().AppendEmpty()
	metric.SetName("network_vm_connections")
	metric.SetUnit("count")
	metric.Metadata().PutStr("instance_id", "06b4020f-461a-4a52-84da-53fa71c2f42b")
	metric.Metadata().PutStr("service.namespace", "SYS.ECS")

	dataPoints = metric.SetEmptyGauge().DataPoints()
	dp = dataPoints.AppendEmpty()
	dp.SetTimestamp(pcommon.NewTimestampFromTime(time.UnixMilli(1556625612000)))
	dp.SetDoubleValue(1)

	dp = dataPoints.AppendEmpty()
	dp.SetTimestamp(pcommon.NewTimestampFromTime(time.UnixMilli(1556625717000)))
	dp.SetDoubleValue(3)
	marshaler := &pmetric.JSONMarshaler{}
	data, err := marshaler.MarshalMetrics(metrics)
	require.NoError(t, err)
	return data
}

func expectedEmptyMetrics(t *testing.T) []byte {
	metrics := pmetric.NewMetrics()
	resourceMetric := metrics.ResourceMetrics().AppendEmpty()

	resource := resourceMetric.Resource()
	resource.Attributes().PutStr("cloud.provider", "huawei_cloud")
	resource.Attributes().PutStr("project.id", "project_1")
	resource.Attributes().PutStr("region", "eu-west-101")

	marshaler := &pmetric.JSONMarshaler{}
	data, err := marshaler.MarshalMetrics(metrics)
	require.NoError(t, err)
	return data
}
