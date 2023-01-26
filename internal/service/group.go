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
		BindNamespace(ctx context.Context, groupId int64, namespace string)
		Unbind(ctx context.Context, groupId int64)
		QueryGroup(ctx context.Context, groupId int64)
		IsGroupBindNamespaceOwnerOrAdmin(ctx context.Context, groupId, userId int64) (yes bool)
		GetApprovalProcess(ctx context.Context, groupId int64) (process map[string]any)
		AddApprovalProcess(ctx context.Context, groupId int64, processName string, args ...string)
		RemoveApprovalProcess(ctx context.Context, groupId int64, processName string, args ...string)
		GetWhitelists(ctx context.Context, groupId int64) (whitelists map[string]any)
		GetBlacklists(ctx context.Context, groupId int64) (blacklists map[string]any)
		GetRegexp(ctx context.Context, groupId int64) (re string)
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
