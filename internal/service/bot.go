// ================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
)

type (
	IBot interface {
		Process(ctx context.Context, ws *ghttp.WebSocket, rawJson []byte, nextProcess func(ctx context.Context))
		CatchEcho(ctx context.Context) (catch bool)
		IsGroupOwnerOrAdmin(ctx context.Context) (yes bool)
		GetPostType(ctx context.Context) string
		GetMsgType(ctx context.Context) string
		GetRequestType(ctx context.Context) string
		GetNoticeType(ctx context.Context) string
		GetSubType(ctx context.Context) string
		GetMsgId(ctx context.Context) int64
		GetMessage(ctx context.Context) string
		GetUserId(ctx context.Context) int64
		GetGroupId(ctx context.Context) int64
		GetComment(ctx context.Context) string
		GetFlag(ctx context.Context) string
		SendMessage(ctx context.Context, messageType string, uid, gid int64, msg string, plain bool)
		SendPlainMsg(ctx context.Context, msg string)
		SendMsg(ctx context.Context, msg string)
		ApproveAddGroup(ctx context.Context, flag, subType string, approve bool, reason string)
		SetModel(ctx context.Context, model string)
		RevokeMessage(ctx context.Context, msgId int64)
		MutePrototype(ctx context.Context, groupId, userId int64, seconds int)
		Mute(ctx context.Context, seconds int)
	}
)

var (
	localBot IBot
)

func Bot() IBot {
	if localBot == nil {
		panic("implement not found for interface IBot, forgot register?")
	}
	return localBot
}

func RegisterBot(i IBot) {
	localBot = i
}
