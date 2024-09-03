// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package huaweicloudreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/huaweicloudreceiver"

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/receiver/receivertest"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
	"go.uber.org/zap/zaptest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/huaweicloudreceiver/internal/mocks"
)

func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}

func TestNewReceiver(t *testing.T) {
	cfg := &Config{
		ControllerConfig: scraperhelper.ControllerConfig{
			CollectionInterval: 1 * time.Second,
		},
	}
	mr := newHuaweiCloudMetricsReceiver(cfg, new(consumertest.MetricsSink), receivertest.NewNopSettings())
	assert.NotNil(t, mr)
}

func TestListMetricDefinitionsSuccess(t *testing.T) {
	mockCes := mocks.NewCesClient(t)

	mockResponse := &model.ListMetricsResponse{
		Metrics: &[]model.MetricInfoList{
			{
				Namespace:  "SYS.ECS",
				MetricName: "cpu_util",
				Dimensions: []model.MetricsDimension{
					{
						Name:  "instance_id",
						Value: "12345",
					},
				},
			},
		},
	}

	mockCes.On("ListMetrics", mock.Anything).Return(mockResponse, nil)

	receiver := &huaweiReceiver{
		clients: &clients{ces: mockCes},
		telType: Metrics,
		config:  createDefaultConfig().(*Config),
	}

	metrics, err := receiver.listMetricDefinitions(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, metrics)
	assert.Equal(t, "SYS.ECS", metrics[0].Namespace)
	assert.Equal(t, "cpu_util", metrics[0].MetricName)
	assert.Equal(t, "instance_id", metrics[0].Dimensions[0].Name)
	assert.Equal(t, "12345", metrics[0].Dimensions[0].Value)
	mockCes.AssertExpectations(t)
}

func TestListMetricDefinitionsFailure(t *testing.T) {
	mockCes := mocks.NewCesClient(t)

	mockCes.On("ListMetrics", mock.Anything).Return(nil, errors.New("failed to list metrics"))
	receiver := &huaweiReceiver{
		clients: &clients{ces: mockCes},
		config:  createDefaultConfig().(*Config),
	}

	metrics, err := receiver.listMetricDefinitions(context.Background())

	assert.Error(t, err)
	assert.Len(t, metrics, 0)
	assert.Equal(t, "failed to list metrics", err.Error())
	mockCes.AssertExpectations(t)
}

func TestListDataPointsForMetricBackOffWIthDefaultConfig(t *testing.T) {
	mockCes := mocks.NewCesClient(t)
	next := new(consumertest.MetricsSink)
	receiver := newHuaweiCloudMetricsReceiver(createDefaultConfig().(*Config), next, receivertest.NewNopSettings())
	receiver.clients = &clients{ces: mockCes}

	mockCes.On("ShowMetricData", mock.Anything).Return(nil, errors.New(requestThrottledErrMsg)).Times(3)
	mockCes.On("ShowMetricData", mock.Anything).Return(&model.ShowMetricDataResponse{
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

	resp, err := receiver.listDataPointsForMetric(context.Background(), time.Now().Add(10*time.Minute), time.Now(), model.MetricInfoList{
		Namespace:  "SYS.ECS",
		MetricName: "cpu_util",
		Dimensions: []model.MetricsDimension{
			{
				Name:  "instance_id",
				Value: "12345",
			},
		},
	})

	require.NoError(t, err)
	assert.Len(t, *resp.Datapoints, 2)
}

func TestListDataPointsForMetricBackOffFails(t *testing.T) {
	mockCes := mocks.NewCesClient(t)
	next := new(consumertest.MetricsSink)
	receiver := newHuaweiCloudReceiver(receivertest.NewNopSettings(), Metrics, &Config{BackOffConfig: configretry.BackOffConfig{
		Enabled:             true,
		InitialInterval:     100 * time.Millisecond,
		MaxInterval:         800 * time.Millisecond,
		MaxElapsedTime:      1 * time.Second,
		RandomizationFactor: 0,
		Multiplier:          2,
	}}, &metricsReceiver{consumer: next})
	receiver.clients = &clients{ces: mockCes}

	mockCes.On("ShowMetricData", mock.Anything).Return(nil, errors.New(requestThrottledErrMsg)).Times(4)

	resp, err := receiver.listDataPointsForMetric(context.Background(), time.Now().Add(10*time.Minute), time.Now(), model.MetricInfoList{
		Namespace:  "SYS.ECS",
		MetricName: "cpu_util",
		Dimensions: []model.MetricsDimension{
			{
				Name:  "instance_id",
				Value: "12345",
			},
		},
	})

	require.ErrorContains(t, err, requestThrottledErrMsg)
	assert.Nil(t, resp)
}

func TestPollMetricsAndConsumeSuccess(t *testing.T) {
	mockCes := mocks.NewCesClient(t)
	next := new(consumertest.MetricsSink)
	receiver := newHuaweiCloudReceiver(receivertest.NewNopSettings(), Metrics, &Config{}, &metricsReceiver{consumer: next})
	receiver.clients = &clients{ces: mockCes}

	mockCes.On("ListMetrics", mock.Anything).Return(&model.ListMetricsResponse{
		Metrics: &[]model.MetricInfoList{
			{
				Namespace:  "SYS.ECS",
				MetricName: "cpu_util",
				Dimensions: []model.MetricsDimension{
					{
						Name:  "instance_id",
						Value: "12345",
					},
				},
			},
		},
	}, nil)

	mockCes.On("ShowMetricData", mock.Anything).Return(&model.ShowMetricDataResponse{
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

	err := receiver.pollTelemetryAndConsume(context.Background())

	require.NoError(t, err)
	assert.Equal(t, 2, next.DataPointCount())
}

func TestStartReadingTelemetry(t *testing.T) {
	tests := []struct {
		name                    string
		scrapeInterval          time.Duration
		setupMocks              func(*mocks.CesClient)
		expectedNumOfDataPoints int
	}{
		{
			name:           "Success case with valid scrape interval",
			scrapeInterval: 2 * time.Second,
			setupMocks: func(m *mocks.CesClient) {
				m.On("ListMetrics", mock.Anything).Return(&model.ListMetricsResponse{
					Metrics: &[]model.MetricInfoList{
						{
							Namespace:  "SYS.ECS",
							MetricName: "cpu_util",
							Dimensions: []model.MetricsDimension{
								{
									Name:  "instance_id",
									Value: "12345",
								},
							},
						},
					},
				}, nil)

				m.On("ShowMetricData", mock.Anything).Return(&model.ShowMetricDataResponse{
					MetricName: stringPtr("cpu_util"),
					Datapoints: &[]model.Datapoint{
						{
							Average:   float64Ptr(45.67),
							Timestamp: 1556625610000,
						},
					},
				}, nil)
			},
			expectedNumOfDataPoints: 1,
		},
		{
			name:           "Error case with Scrape returning error",
			scrapeInterval: 1 * time.Second,
			setupMocks: func(m *mocks.CesClient) {
				m.On("ListMetrics", mock.Anything).Return(nil, errors.New("server error"))
			},
			expectedNumOfDataPoints: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCes := mocks.NewCesClient(t)
			next := new(consumertest.MetricsSink)
			tt.setupMocks(mockCes)
			logger := zaptest.NewLogger(t)
			r := &huaweiReceiver{
				config: &Config{
					ControllerConfig: scraperhelper.ControllerConfig{
						CollectionInterval: tt.scrapeInterval,
						InitialDelay:       10 * time.Millisecond,
					},
				},
				clients:       &clients{ces: mockCes},
				telType:       Metrics,
				logger:        logger,
				dataProcessor: &metricsReceiver{consumer: next},
				lastSeenTs:    make(map[string]time.Time),
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r.startReadingTelemetry(ctx)

			assert.Equal(t, tt.expectedNumOfDataPoints, next.DataPointCount())
		})
	}
}
func TestCreateHTTPConfigNoVerifySSL(t *testing.T) {
	cfg, err := createHTTPConfig(HuaweiSessionConfig{NoVerifySSL: true})
	require.NoError(t, err)
	assert.Equal(t, cfg.IgnoreSSLVerification, true)
}

func TestCreateHTTPConfigWithProxy(t *testing.T) {
	cfg, err := createHTTPConfig(HuaweiSessionConfig{
		ProxyAddress:  "https://127.0.0.1:8888",
		ProxyUser:     "admin",
		ProxyPassword: "pass",
		AccessKey:     "123",
		SecretKey:     "secret",
	})
	require.NoError(t, err)
	assert.Equal(t, cfg.HttpProxy.Schema, "https")
	assert.Equal(t, cfg.HttpProxy.Host, "127.0.0.1")
	assert.Equal(t, cfg.HttpProxy.Port, 8888)
	assert.Equal(t, cfg.IgnoreSSLVerification, false)

}
