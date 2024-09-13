// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/huaweicloudlogsreceiver/internal"

import (
	"fmt"
	"regexp"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/lts/v2/model"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
)

const (
	layout = "2006-01-02/15:04:05"
	// Define the regex pattern to match the timestamp
	pattern = `(\d{4}-\d{2}-\d{2}/\d{2}:\d{2}:\d{2})`
)

var (
	re = regexp.MustCompile(pattern)
)

// example of log content : "2020-07-25/14:44:43 this <HighLightTag>log</HighLightTag> is Error NO 2"
func extractTimestamp(log string) (time.Time, error) {
	matches := re.FindStringSubmatch(log)
	if len(matches) == 0 {
		return time.Time{}, fmt.Errorf("timestamp not found")
	}

	// Note: Assumtion is that the date string is using UTC timezone.
	timestamp, err := time.Parse(layout, matches[0])
	if err != nil {
		return time.Time{}, err
	}
	return timestamp, nil
}

func ConvertLTSLogsToOTLP(projectID, regionID, groupID, streamD string, ltsLogs []model.LogContents) plog.Logs {
	logs := plog.NewLogs()
	resourceLog := logs.ResourceLogs().AppendEmpty()

	resource := resourceLog.Resource()
	resourceAttr := resource.Attributes()
	resourceAttr.PutStr("cloud.provider", "huawei_cloud")
	resourceAttr.PutStr("project.id", projectID)
	resourceAttr.PutStr("region.id", regionID)
	resourceAttr.PutStr("group.id", groupID)
	resourceAttr.PutStr("stream.id", streamD)

	if len(ltsLogs) == 0 {
		return logs
	}
	scopedLog := resourceLog.ScopeLogs().AppendEmpty()
	scopedLog.Scope().SetName("huawei_cloud_lts")
	scopedLog.Scope().SetVersion("v2")
	for _, ltsLog := range ltsLogs {

		logRecord := scopedLog.LogRecords().AppendEmpty()
		if ltsLog.Content != nil {
			content := *ltsLog.Content
			if ts, err := extractTimestamp(content); err == nil {
				logRecord.SetObservedTimestamp(pcommon.Timestamp(ts.UnixNano()))
			}
			logRecord.Body().SetStr(content)
		}
		for key, value := range ltsLog.Labels {
			logRecord.Attributes().PutStr(key, value)
		}
		if ltsLog.LineNum != nil {
			logRecord.Attributes().PutStr("lineNum", *ltsLog.LineNum)
		}
	}

	return logs
}
