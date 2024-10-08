// Code generated by mockery v2.44.1. DO NOT EDIT.

package mocks

import (
	model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"
	mock "github.com/stretchr/testify/mock"
)

// CesClient is an autogenerated mock type for the CesClient type
type CesClient struct {
	mock.Mock
}

// ListMetrics provides a mock function with given fields: request
func (_m *CesClient) ListMetrics(request *model.ListMetricsRequest) (*model.ListMetricsResponse, error) {
	ret := _m.Called(request)

	if len(ret) == 0 {
		panic("no return value specified for ListMetrics")
	}

	var r0 *model.ListMetricsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.ListMetricsRequest) (*model.ListMetricsResponse, error)); ok {
		return rf(request)
	}
	if rf, ok := ret.Get(0).(func(*model.ListMetricsRequest) *model.ListMetricsResponse); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ListMetricsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.ListMetricsRequest) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ShowMetricData provides a mock function with given fields: request
func (_m *CesClient) ShowMetricData(request *model.ShowMetricDataRequest) (*model.ShowMetricDataResponse, error) {
	ret := _m.Called(request)

	if len(ret) == 0 {
		panic("no return value specified for ShowMetricData")
	}

	var r0 *model.ShowMetricDataResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.ShowMetricDataRequest) (*model.ShowMetricDataResponse, error)); ok {
		return rf(request)
	}
	if rf, ok := ret.Get(0).(func(*model.ShowMetricDataRequest) *model.ShowMetricDataResponse); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ShowMetricDataResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.ShowMetricDataRequest) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewCesClient creates a new instance of CesClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCesClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *CesClient {
	mock := &CesClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
