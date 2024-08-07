// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"github.com/bytedance/sonic/ast"
	"github.com/gorilla/websocket"
)

type (
	IBot interface {
		CtxWithWebSocket(parent context.Context, conn *websocket.Conn) context.Context
		CtxNewWebSocketMutex(parent context.Context) context.Context
		CtxWithReqJson(ctx context.Context, reqJson *ast.Node) context.Context
		Process(ctx context.Context, rawJson []byte, nextProcess func(ctx context.Context))
		JoinConnectionPool(ctx context.Context, key int64)
		LeaveConnectionPool(key int64)
		LoadConnectionPool(key int64) context.Context
		Forward(ctx context.Context, url string, authorization string) error
		GetPostType(ctx context.Context) string
		GetMsgType(ctx context.Context) string
		GuessMsgType(groupId int64) string
		GetRequestType(ctx context.Context) string
		GetNoticeType(ctx context.Context) string
		GetSubType(ctx context.Context) string
		GetMsgId(ctx context.Context) int64
		GetMessage(ctx context.Context) string
		GetUserId(ctx context.Context) int64
		GetGroupId(ctx context.Context) int64
		GetComment(ctx context.Context) string
		GetFlag(ctx context.Context) string
		GetTimestamp(ctx context.Context) int64
		GetOperatorId(ctx context.Context) int64
		GetSelfId(ctx context.Context) int64
		GetNickname(ctx context.Context) string
		GetCard(ctx context.Context) string
		GetCardOrNickname(ctx context.Context) string
		GetCardOldNew(ctx context.Context) (oldCard string, newCard string)
		GetGroupMemberInfo(ctx context.Context, groupId int64, userId int64) (member ast.Node, err error)
		GetGroupMemberList(ctx context.Context, groupId int64, usingCache ...bool) (members []any, err error)
		RequestMessage(ctx context.Context, messageId int64) (messageMap map[string]any, err error)
		GetGroupInfo(ctx context.Context, groupId int64, noCache ...bool) (infoMap map[string]any, err error)
		GetLoginInfo(ctx context.Context) (userId int64, nickname string)
		IsGroupOwnerOrAdmin(ctx context.Context) bool
		IsGroupOwnerOrAdminOrSysTrusted(ctx context.Context) bool
		SendMessage(ctx context.Context, messageType string, userId int64, groupId int64, msg string, plain bool) (err error)
		SendPlainMsg(ctx context.Context, msg string)
		SendMsg(ctx context.Context, msg string)
		SendPlainMsgIfNotApiReq(ctx context.Context, msg string)
		SendFileToGroup(ctx context.Context, groupId int64, filePath string, name string, folder string)
		SendFileToUser(ctx context.Context, userId int64, filePath string, name string)
		SendFile(ctx context.Context, filePath string, name string)
		UploadFile(ctx context.Context, url string) (filePath string, err error)
		ApproveJoinGroup(ctx context.Context, flag string, subType string, approve bool, reason string)
		SetModel(ctx context.Context, model string)
		RecallMessage(ctx context.Context, msgId int64)
		MutePrototype(ctx context.Context, groupId int64, userId int64, seconds int)
		Mute(ctx context.Context, seconds int)
		SetGroupCard(ctx context.Context, groupId int64, userId int64, card string)
		Kick(ctx context.Context, groupId int64, userId int64, reject ...bool)
		RewriteMessage(ctx context.Context, message string)
		SetHistory(ctx context.Context, history string) error
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
