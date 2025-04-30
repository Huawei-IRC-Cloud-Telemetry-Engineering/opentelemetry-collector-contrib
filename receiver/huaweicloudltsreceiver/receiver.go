// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package huaweicloudltsreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/huaweicloudltsreceiver"

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	lts "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/lts/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/lts/v2/model"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/lts/v2/region"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/huawei"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/huaweicloudltsreceiver/internal"
)

const (
	// See https://support.huaweicloud.com/intl/en-us/devg-apisign/api-sign-errorcode.html
	requestThrottledErrMsg = "APIGW.0308"
)

type logsReceiver struct {
	logger *zap.Logger
	client internal.LtsClient
	cancel context.CancelFunc

	host         component.Host
	nextConsumer consumer.Logs
	lastTs       time.Time
	config       *Config
	shutdownChan chan struct{}
}

func newhuaweicloudltsreceiver(settings receiver.Settings, cfg *Config, next consumer.Logs) *logsReceiver {
	rcvr := &logsReceiver{
		logger:       settings.Logger,
		config:       cfg,
		nextConsumer: next,
		shutdownChan: make(chan struct{}, 1),
	}
	return rcvr
}

func (rcvr *logsReceiver) Start(ctx context.Context, host component.Host) error {
	rcvr.host = host
	ctx, rcvr.cancel = context.WithCancel(ctx)

	if rcvr.client == nil {
		client, err := rcvr.createClient()
		if err != nil {
			rcvr.logger.Error(err.Error())
			return nil
		}
		rcvr.client = client
	}

	go rcvr.startReadingLogs(ctx)
	return nil
}

func (rcvr *logsReceiver) startReadingLogs(ctx context.Context) {
	if rcvr.config.InitialDelay > 0 {
		<-time.After(rcvr.config.InitialDelay)
	}
	rcvr.logger.Info("Starting to poll logs from Huawei Cloud LTS", zap.String("log_group_id", rcvr.config.GroupID), zap.String("log_stream_id", rcvr.config.StreamID))

	if err := rcvr.pollLogsAndConsume(ctx); err != nil {
		rcvr.logger.Error(err.Error())
	}
	ticker := time.NewTicker(rcvr.config.CollectionInterval)

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			//  TODO: Improve error handling for client-server interactions
			//  The current implementation lacks robust error handling, especially for
			//  scenarios such as service unavailability, timeouts, and request errors.
			//  - Investigate how to handle service unavailability or timeouts gracefully.
			//  - Implement appropriate actions or retries for different types of request errors.
			//  - Refer to the Huawei SDK documentation to identify
			//    all possible client/request errors and determine how to manage them.
			//  - Consider implementing custom error messages or fallback mechanisms for critical failures.

			if err := rcvr.pollLogsAndConsume(ctx); err != nil {
				rcvr.logger.Error(err.Error())
			}
		case <-ctx.Done():
			return
		}
	}
}

func (rcvr *logsReceiver) createClient() (*lts.LtsClient, error) {
	auth, err := basic.NewCredentialsBuilder().
		WithAk(string(rcvr.config.AccessKey)).
		WithSk(string(rcvr.config.SecretKey)).
		WithProjectId(rcvr.config.ProjectID).
		SafeBuild()

	if err != nil {
		return nil, err
	}

	httpConfig, err := huawei.CreateHTTPConfig(rcvr.config.HuaweiSessionConfig)
	if err != nil {
		return nil, err
	}
	r, err := region.SafeValueOf(rcvr.config.RegionID)
	if err != nil {
		return nil, err
	}

	hcHTTPConfig, err := lts.LtsClientBuilder().
		WithRegion(r).
		WithCredential(auth).
		WithHttpConfig(httpConfig).
		SafeBuild()

	if err != nil {
		return nil, err
	}

	return lts.NewLtsClient(hcHTTPConfig), nil
}

func (rcvr *logsReceiver) pollLogsAndConsume(ctx context.Context) error {
	if rcvr.client == nil {
		return errors.New("invalid client")
	}
	to := time.Now()
	from := rcvr.lastTs
	if from.IsZero() {
		from = to.Add(-1 * rcvr.config.CollectionInterval)
	}
	logs, err := rcvr.listLogs(ctx, from, to)
	if err != nil {
		return err
	}

	if len(logs) == 0 {
		rcvr.logger.Info("No logs found in the specified time range", zap.Time("from", from), zap.Time("to", to))
		return nil
	}

	config := rcvr.config
	optLogs := internal.ConvertLTSLogsToOTLP(config.ProjectID, config.RegionID, config.GroupID, config.StreamID, logs)
	if err := rcvr.nextConsumer.ConsumeLogs(ctx, optLogs); err != nil {
		return err
	}
	rcvr.lastTs = to
	return nil
}

func (rcvr *logsReceiver) listLogs(ctx context.Context, from, to time.Time) ([]model.LogContents, error) {
	// TODO: Add pagination logic. Check IsQueryComplete field
	response, err := huawei.MakeAPICallWithRetry(
		ctx,
		rcvr.shutdownChan,
		rcvr.logger,
		func() (*model.ListLogsResponse, error) {
			return rcvr.client.ListLogs(&model.ListLogsRequest{
				LogGroupId:  rcvr.config.GroupID,
				LogStreamId: rcvr.config.StreamID,
				Body: &model.QueryLtsLogParams{
					StartTime: strconv.FormatInt(from.UTC().UnixMilli(), 10),
					EndTime:   strconv.FormatInt(to.UTC().UnixMilli(), 10),
				},
			})
		},
		func(err error) bool { return strings.Contains(err.Error(), requestThrottledErrMsg) },
		huawei.NewExponentialBackOff(&rcvr.config.BackOffConfig),
	)
	if err != nil {
		return []model.LogContents{}, err
	}

	if response == nil {
		return nil, errors.New("invalid response from Huawei Cloud LTS API")
	}

	if response.Logs == nil {
		return []model.LogContents{}, nil
	}

	return *response.Logs, nil
}

func (rcvr *logsReceiver) Shutdown(_ context.Context) error {
	if rcvr.cancel != nil {
		rcvr.cancel()
	}
	rcvr.shutdownChan <- struct{}{}
	close(rcvr.shutdownChan)
	return nil
}
