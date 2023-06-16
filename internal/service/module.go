// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"time"
)

type (
	IModule interface {
		TryApproveAddGroup(ctx context.Context) (catch bool)
		TryLockCard(ctx context.Context) (catch bool)
		TryAutoSetCard(ctx context.Context) (catch bool)
		TryKeywordRecall(ctx context.Context) (catch bool)
		TryKeywordReply(ctx context.Context) (catch bool)
		TryLogLeave(ctx context.Context) (catch bool)
		TryUndoMessageRecall(ctx context.Context) (catch bool)
		MultiContains(str string, m map[string]any) (contains bool, hit string, mValue string)
		AutoMute(ctx context.Context, kind string, groupId, userId int64, passTimes, baseMinutes, limitMinutes int, duration time.Duration)
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
