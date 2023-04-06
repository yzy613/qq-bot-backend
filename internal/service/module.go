// ================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	IModule interface {
		TryApproveAddGroup(ctx context.Context) (catch bool)
		TryKeywordRevoke(ctx context.Context) (catch bool)
		TryKeywordReply(ctx context.Context) (catch bool)
		TryLogLeave(ctx context.Context) (catch bool)
		MultiContains(str string, m map[string]any) (contains bool, hit string, mValue string)
	}
)

var (
	localModule IModule
)

func Module() IModule {
	if localModule == nil {
		panic("implement not found for interface IModule, forgot register?")
	}
	return localModule
}

func RegisterModule(i IModule) {
	localModule = i
}
