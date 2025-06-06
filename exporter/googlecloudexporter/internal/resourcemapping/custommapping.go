// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resourcemapping // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudexporter/internal/resourcemapping"

import (
	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/collector"
	"go.opentelemetry.io/collector/pdata/pcommon"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
)

var (
	mappingKey                   = "gcp.resource_type"
	monitoredResourceLabelPrefix = "gcp."
)

// CustomLoggingMonitoredResourceMapping allows mapping from OTel resources to
// Monitored Resources defined here:
// https://cloud.google.com/monitoring/api/resources
//
// The monitored resource type is extracted from the `gcp.resource_type`
// attribute. And the monitored resource labels are extracted from resource
// attributes with the prefix `gcp.<monitored resource type>.`.
func CustomLoggingMonitoredResourceMapping(r pcommon.Resource) *monitoredrespb.MonitoredResource {
	return customMonitoredResourceMapping(r, collector.DefaultConfig().LogConfig.MapMonitoredResource)
}

// CustomMetricMonitoredResourceMapping allows mapping from OTel resources to
// Monitored Resources defined here:
// https://cloud.google.com/monitoring/api/resources
//
// The monitored resource type is extracted from the `gcp.resource_type`
// attribute. And the monitored resource labels are extracted from resource
// attributes with the prefix `gcp.<monitored resource type>.`.
func CustomMetricMonitoredResourceMapping(r pcommon.Resource) *monitoredrespb.MonitoredResource {
	return customMonitoredResourceMapping(r, collector.DefaultConfig().MetricConfig.MapMonitoredResource)
}

func customMonitoredResourceMapping(r pcommon.Resource, mmrFunc func(pcommon.Resource) *monitoredrespb.MonitoredResource) *monitoredrespb.MonitoredResource {
	var monitoredResourceType string
	monitoredResourceLabels := make(map[string]string)
	for k, v := range r.Attributes().All() {
		if k == mappingKey {
			monitoredResourceType = v.AsString()
			break
		}
	}

	if monitoredResourceType == "" {
		return mmrFunc(r)
	}

	prefix := monitoredResourceLabelPrefix + monitoredResourceType + "."
	for k, v := range r.Attributes().All() {
		// Extract the monitored resource label by separating it from the prefix.
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			monitoredResourceLabels[k[len(prefix):]] = v.AsString()
		}
	}

	return &monitoredrespb.MonitoredResource{
		Type:   monitoredResourceType,
		Labels: monitoredResourceLabels,
	}
}
