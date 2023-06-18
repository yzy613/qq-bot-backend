// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	IGroup interface {
		GetApprovalProcess(ctx context.Context, groupId int64) (process map[string]any)
		GetApprovalWhitelists(ctx context.Context, groupId int64) (whitelists map[string]any)
		GetApprovalBlacklists(ctx context.Context, groupId int64) (blacklists map[string]any)
		GetApprovalRegexp(ctx context.Context, groupId int64) (exp string)
		IsEnabledApprovalAutoPass(ctx context.Context, groupId int64) (enabled bool)
		AddApprovalProcessWithRes(ctx context.Context, groupId int64, processName string, args ...string)
		RemoveApprovalProcessWithRes(ctx context.Context, groupId int64, processName string, args ...string)
		GetCardAutoSetList(ctx context.Context, groupId int64) (listName string)
		IsCardLocked(ctx context.Context, groupId int64) (locked bool)
		SetAutoSetListWithRes(ctx context.Context, groupId int64, listName string)
		RemoveAutoSetListWithRes(ctx context.Context, groupId int64)
		CheckCardWithRegexpWithRes(ctx context.Context, groupId int64, listName, exp string)
		CheckCardByListWithRes(ctx context.Context, groupId int64, toList, fromList string)
		LockCardWithRes(ctx context.Context, groupId int64)
		UnlockCardWithRes(ctx context.Context, groupId int64)
		ExportGroupMemberListWithRes(ctx context.Context, groupId int64, listName string)
		BindNamespaceWithRes(ctx context.Context, groupId int64, namespace string)
		UnbindWithRes(ctx context.Context, groupId int64)
		QueryGroupWithRes(ctx context.Context, groupId int64)
		KickFromListWithRes(ctx context.Context, groupId int64, listName string)
		KeepFromListWithRes(ctx context.Context, groupId int64, listName string)
		GetKeywordProcess(ctx context.Context, groupId int64) (process map[string]any)
		GetKeywordWhitelists(ctx context.Context, groupId int64) (whitelists map[string]any)
		GetKeywordBlacklists(ctx context.Context, groupId int64) (blacklists map[string]any)
		GetKeywordReplyList(ctx context.Context, groupId int64) (listName string)
		AddKeywordProcessWithRes(ctx context.Context, groupId int64, processName string, args ...string)
		RemoveKeywordProcessWithRes(ctx context.Context, groupId int64, processName string, args ...string)
		GetLogLeaveList(ctx context.Context, groupId int64) (listName string)
		SetLogLeaveListWithRes(ctx context.Context, groupId int64, listName string)
		RemoveLogLeaveListWithRes(ctx context.Context, groupId int64)
		IsEnabledAntiRecall(ctx context.Context, groupId int64) (enabled bool)
		SetAntiRecallWithRes(ctx context.Context, groupId int64, enable bool)
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
