// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package networkscraper // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver/internal/scraper/networkscraper"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver/internal"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver/internal/scraper/networkscraper/internal/metadata"
)

// This file implements Factory for Network scraper.

const (
	// TypeStr the value of "type" key in configuration.
	TypeStr = "network"
)

var (
	// scraperType is the component type used for the built scraper.
	scraperType component.Type = component.MustNewType(TypeStr)
)

// Factory is the Factory for scraper.
type Factory struct {
}

// CreateDefaultConfig creates the default configuration for the Scraper.
func (f *Factory) CreateDefaultConfig() internal.Config {
	return &Config{
		MetricsBuilderConfig: metadata.DefaultMetricsBuilderConfig(),
	}
}

// CreateMetricsScraper creates a scraper based on provided config.
func (f *Factory) CreateMetricsScraper(
	ctx context.Context,
	settings receiver.Settings,
	config internal.Config,
) (scraperhelper.Scraper, error) {
	cfg := config.(*Config)
	s, err := newNetworkScraper(ctx, settings, cfg)
	if err != nil {
		return nil, err
	}

	return scraperhelper.NewScraper(
		scraperType,
		s.scrape,
		scraperhelper.WithStart(s.start),
	)
}
