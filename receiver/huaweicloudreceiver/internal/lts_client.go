// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/huaweicloudreceiver/internal"

import (
	lts "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/lts/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/lts/v2/model"
)

// This interface should have all the function defined inside https://github.com/huaweicloud/huaweicloud-sdk-go-v3/blob/v0.1.110/services/lts/v2/lts_client.go
// Check https://github.com/vektra/mockery on how to install it on your machine.
//
//go:generate mockery --name LtsClient --case=underscore --output=./mocks
type LtsClient interface {
	CreateAccessConfig(request *model.CreateAccessConfigRequest) (*model.CreateAccessConfigResponse, error)
	CreateAccessConfigInvoker(request *model.CreateAccessConfigRequest) *lts.CreateAccessConfigInvoker

	DeleteAccessConfig(request *model.DeleteAccessConfigRequest) (*model.DeleteAccessConfigResponse, error)
	DeleteAccessConfigInvoker(request *model.DeleteAccessConfigRequest) *lts.DeleteAccessConfigInvoker

	ListAccessConfig(request *model.ListAccessConfigRequest) (*model.ListAccessConfigResponse, error)
	ListAccessConfigInvoker(request *model.ListAccessConfigRequest) *lts.ListAccessConfigInvoker

	UpdateAccessConfig(request *model.UpdateAccessConfigRequest) (*model.UpdateAccessConfigResponse, error)
	UpdateAccessConfigInvoker(request *model.UpdateAccessConfigRequest) *lts.UpdateAccessConfigInvoker

	CreateHostGroup(request *model.CreateHostGroupRequest) (*model.CreateHostGroupResponse, error)
	CreateHostGroupInvoker(request *model.CreateHostGroupRequest) *lts.CreateHostGroupInvoker

	DeleteHostGroup(request *model.DeleteHostGroupRequest) (*model.DeleteHostGroupResponse, error)
	DeleteHostGroupInvoker(request *model.DeleteHostGroupRequest) *lts.DeleteHostGroupInvoker

	ListHostGroup(request *model.ListHostGroupRequest) (*model.ListHostGroupResponse, error)
	ListHostGroupInvoker(request *model.ListHostGroupRequest) *lts.ListHostGroupInvoker

	UpdateHostGroup(request *model.UpdateHostGroupRequest) (*model.UpdateHostGroupResponse, error)
	UpdateHostGroupInvoker(request *model.UpdateHostGroupRequest) *lts.UpdateHostGroupInvoker

	CreateKeywordsAlarmRule(request *model.CreateKeywordsAlarmRuleRequest) (*model.CreateKeywordsAlarmRuleResponse, error)
	CreateKeywordsAlarmRuleInvoker(request *model.CreateKeywordsAlarmRuleRequest) *lts.CreateKeywordsAlarmRuleInvoker

	DeleteKeywordsAlarmRule(request *model.DeleteKeywordsAlarmRuleRequest) (*model.DeleteKeywordsAlarmRuleResponse, error)
	DeleteKeywordsAlarmRuleInvoker(request *model.DeleteKeywordsAlarmRuleRequest) *lts.DeleteKeywordsAlarmRuleInvoker

	ListKeywordsAlarmRules(request *model.ListKeywordsAlarmRulesRequest) (*model.ListKeywordsAlarmRulesResponse, error)
	ListKeywordsAlarmRulesInvoker(request *model.ListKeywordsAlarmRulesRequest) *lts.ListKeywordsAlarmRulesInvoker

	UpdateKeywordsAlarmRule(request *model.UpdateKeywordsAlarmRuleRequest) (*model.UpdateKeywordsAlarmRuleResponse, error)
	UpdateKeywordsAlarmRuleInvoker(request *model.UpdateKeywordsAlarmRuleRequest) *lts.UpdateKeywordsAlarmRuleInvoker

	CreateLogDumpObs(request *model.CreateLogDumpObsRequest) (*model.CreateLogDumpObsResponse, error)
	CreateLogDumpObsInvoker(request *model.CreateLogDumpObsRequest) *lts.CreateLogDumpObsInvoker

	CreateLogGroup(request *model.CreateLogGroupRequest) (*model.CreateLogGroupResponse, error)
	CreateLogGroupInvoker(request *model.CreateLogGroupRequest) *lts.CreateLogGroupInvoker

	DeleteLogGroup(request *model.DeleteLogGroupRequest) (*model.DeleteLogGroupResponse, error)
	DeleteLogGroupInvoker(request *model.DeleteLogGroupRequest) *lts.DeleteLogGroupInvoker

	ListLogGroups(request *model.ListLogGroupsRequest) (*model.ListLogGroupsResponse, error)
	ListLogGroupsInvoker(request *model.ListLogGroupsRequest) *lts.ListLogGroupsInvoker

	UpdateLogGroup(request *model.UpdateLogGroupRequest) (*model.UpdateLogGroupResponse, error)
	UpdateLogGroupInvoker(request *model.UpdateLogGroupRequest) *lts.UpdateLogGroupInvoker

	DeleteLogStream(request *model.DeleteLogStreamRequest) (*model.DeleteLogStreamResponse, error)
	DeleteLogStreamInvoker(request *model.DeleteLogStreamRequest) *lts.DeleteLogStreamInvoker

	ListLogStreams(request *model.ListLogStreamsRequest) (*model.ListLogStreamsResponse, error)
	ListLogStreamsInvoker(request *model.ListLogStreamsRequest) *lts.ListLogStreamsInvoker

	UpdateLogStream(request *model.UpdateLogStreamRequest) (*model.UpdateLogStreamResponse, error)
	UpdateLogStreamInvoker(request *model.UpdateLogStreamRequest) *lts.UpdateLogStreamInvoker

	CreateNotificationTemplate(request *model.CreateNotificationTemplateRequest) (*model.CreateNotificationTemplateResponse, error)
	CreateNotificationTemplateInvoker(request *model.CreateNotificationTemplateRequest) *lts.CreateNotificationTemplateInvoker

	DeleteNotificationTemplate(request *model.DeleteNotificationTemplateRequest) (*model.DeleteNotificationTemplateResponse, error)
	DeleteNotificationTemplateInvoker(request *model.DeleteNotificationTemplateRequest) *lts.DeleteNotificationTemplateInvoker

	ListNotificationTemplates(request *model.ListNotificationTemplatesRequest) (*model.ListNotificationTemplatesResponse, error)
	ListNotificationTemplatesInvoker(request *model.ListNotificationTemplatesRequest) *lts.ListNotificationTemplatesInvoker

	ShowNotificationTemplate(request *model.ShowNotificationTemplateRequest) (*model.ShowNotificationTemplateResponse, error)
	ShowNotificationTemplateInvoker(request *model.ShowNotificationTemplateRequest) *lts.ShowNotificationTemplateInvoker

	UpdateNotificationTemplate(request *model.UpdateNotificationTemplateRequest) (*model.UpdateNotificationTemplateResponse, error)
	UpdateNotificationTemplateInvoker(request *model.UpdateNotificationTemplateRequest) *lts.UpdateNotificationTemplateInvoker

	CreateSearchCriterias(request *model.CreateSearchCriteriasRequest) (*model.CreateSearchCriteriasResponse, error)
	CreateSearchCriteriasInvoker(request *model.CreateSearchCriteriasRequest) *lts.CreateSearchCriteriasInvoker

	DeleteSearchCriterias(request *model.DeleteSearchCriteriasRequest) (*model.DeleteSearchCriteriasResponse, error)
	DeleteSearchCriteriasInvoker(request *model.DeleteSearchCriteriasRequest) *lts.DeleteSearchCriteriasInvoker

	ListCriterias(request *model.ListCriteriasRequest) (*model.ListCriteriasResponse, error)
	ListCriteriasInvoker(request *model.ListCriteriasRequest) *lts.ListCriteriasInvoker

	CreateStructConfig(request *model.CreateStructConfigRequest) (*model.CreateStructConfigResponse, error)
	CreateStructConfigInvoker(request *model.CreateStructConfigRequest) *lts.CreateStructConfigInvoker

	ListStructTemplate(request *model.ListStructTemplateRequest) (*model.ListStructTemplateResponse, error)
	ListStructTemplateInvoker(request *model.ListStructTemplateRequest) *lts.ListStructTemplateInvoker

	UpdateStructConfig(request *model.UpdateStructConfigRequest) (*model.UpdateStructConfigResponse, error)
	UpdateStructConfigInvoker(request *model.UpdateStructConfigRequest) *lts.UpdateStructConfigInvoker

	CreateTransfer(request *model.CreateTransferRequest) (*model.CreateTransferResponse, error)
	CreateTransferInvoker(request *model.CreateTransferRequest) *lts.CreateTransferInvoker

	DeleteTransfer(request *model.DeleteTransferRequest) (*model.DeleteTransferResponse, error)
	DeleteTransferInvoker(request *model.DeleteTransferRequest) *lts.DeleteTransferInvoker

	ListTransfers(request *model.ListTransfersRequest) (*model.ListTransfersResponse, error)
	ListTransfersInvoker(request *model.ListTransfersRequest) *lts.ListTransfersInvoker

	UpdateTransfer(request *model.UpdateTransferRequest) (*model.UpdateTransferResponse, error)
	UpdateTransferInvoker(request *model.UpdateTransferRequest) *lts.UpdateTransferInvoker

	Createfavorite(request *model.CreatefavoriteRequest) (*model.CreatefavoriteResponse, error)
	CreatefavoriteInvoker(request *model.CreatefavoriteRequest) *lts.CreatefavoriteInvoker

	Deletefavorite(request *model.DeletefavoriteRequest) (*model.DeletefavoriteResponse, error)
	DeletefavoriteInvoker(request *model.DeletefavoriteRequest) *lts.DeletefavoriteInvoker

	CreateSqlAlarmRule(request *model.CreateSqlAlarmRuleRequest) (*model.CreateSqlAlarmRuleResponse, error)
	CreateSqlAlarmRuleInvoker(request *model.CreateSqlAlarmRuleRequest) *lts.CreateSqlAlarmRuleInvoker

	DeleteSqlAlarmRule(request *model.DeleteSqlAlarmRuleRequest) (*model.DeleteSqlAlarmRuleResponse, error)
	DeleteSqlAlarmRuleInvoker(request *model.DeleteSqlAlarmRuleRequest) *lts.DeleteSqlAlarmRuleInvoker

	ListSqlAlarmRules(request *model.ListSqlAlarmRulesRequest) (*model.ListSqlAlarmRulesResponse, error)
	ListSqlAlarmRulesInvoker(request *model.ListSqlAlarmRulesRequest) *lts.ListSqlAlarmRulesInvoker

	UpdateSqlAlarmRule(request *model.UpdateSqlAlarmRuleRequest) (*model.UpdateSqlAlarmRuleResponse, error)
	UpdateSqlAlarmRuleInvoker(request *model.UpdateSqlAlarmRuleRequest) *lts.UpdateSqlAlarmRuleInvoker

	CreateLogStream(request *model.CreateLogStreamRequest) (*model.CreateLogStreamResponse, error)
	CreateLogStreamInvoker(request *model.CreateLogStreamRequest) *lts.CreateLogStreamInvoker
}
