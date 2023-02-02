// ================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	IGroup interface {
		AddApprovalProcess(ctx context.Context, groupId int64, processName string, args ...string)
		RemoveApprovalProcess(ctx context.Context, groupId int64, processName string, args ...string)
		GetApprovalProcess(ctx context.Context, groupId int64) (process map[string]any)
		GetApprovalWhitelists(ctx context.Context, groupId int64) (whitelists map[string]any)
		GetApprovalBlacklists(ctx context.Context, groupId int64) (blacklists map[string]any)
		GetApprovalRegexp(ctx context.Context, groupId int64) (exp string)
		BindNamespace(ctx context.Context, groupId int64, namespace string)
		Unbind(ctx context.Context, groupId int64)
		QueryGroup(ctx context.Context, groupId int64)
		ExportGroupMemberList(ctx context.Context, groupId int64, listName string)
		AddKeywordProcess(ctx context.Context, groupId int64, processName string, args ...string)
		RemoveKeywordProcess(ctx context.Context, groupId int64, processName string, args ...string)
		GetKeywordProcess(ctx context.Context, groupId int64) (process map[string]any)
		GetKeywordWhitelists(ctx context.Context, groupId int64) (whitelists map[string]any)
		GetKeywordBlacklists(ctx context.Context, groupId int64) (blacklists map[string]any)
		SetLogLeaveList(ctx context.Context, groupId int64, listName string)
		RemoveLogLeaveList(ctx context.Context, groupId int64)
		GetLogLeaveList(ctx context.Context, groupId int64) (listName string)
	}
)

var (
	localGroup IGroup
)

func Group() IGroup {
	if localGroup == nil {
		panic("implement not found for interface IGroup, forgot register?")
	}
	return localGroup
}

func RegisterGroup(i IGroup) {
	localGroup = i
}
