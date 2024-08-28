package event

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
	"regexp"
)

var (
	cqReplyRe = regexp.MustCompile(`\[CQ:reply,.+?]`)
)

func (s *sEvent) TryUndoMessageRecall(ctx context.Context) (catch bool) {
	{
		userId := service.Bot().GetUserId(ctx)
		// 不处理自己的消息
		if service.Bot().GetSelfId(ctx) == userId {
			return
		}
		// 仅处理发送者自己撤回的消息
		if service.Bot().GetOperatorId(ctx) != userId {
			return
		}
		groupId := service.Bot().GetGroupId(ctx)
		// owner or admin 在 only-anti-recall-member 情况下不需要反撤回
		if service.Group().IsOnlyAntiRecallMemberSet(ctx, groupId) && service.Bot().IsGroupOwnerOrAdmin(ctx) {
			return
		}
		// 获取当前 group message anti-recall 策略
		if !service.Group().IsAntiRecallEnabled(ctx, groupId) {
			return
		}
	}

	// 获取撤回消息的 id
	msgId := service.Bot().GetMsgId(ctx)
	// 获取消息
	messageMap, err := service.Bot().RequestMessage(ctx, msgId)
	if err != nil {
		service.Bot().SendPlainMsg(ctx, "获取历史消息失败")
		return
	}
	// 获取消息信息
	senderMap := gconv.Map(messageMap["sender"])
	nickname := gconv.String(senderMap["card"])
	if nickname == "" {
		nickname = gconv.String(senderMap["nickname"])
	}
	userId := gconv.Int64(senderMap["user_id"])
	groupId := gconv.Int64(messageMap["group_id"])
	// 获取撤回的消息
	message := gconv.String(messageMap["message"])
	// 防止过度触发反撤回
	service.Util().AutoMute(ctx, "recall", groupId, userId,
		2, 5, 5, gconv.Duration("1m"))

	// 反撤回
	notificationGroupId := service.Group().GetMessageNotificationGroupId(ctx, groupId)
	var msg string
	if notificationGroupId == 0 {
		notificationGroupId = groupId
		msg = nickname + "(" + gconv.String(userId) + ") 撤回了：\n"
	} else {
		msg = nickname + "(" + gconv.String(userId) +
			") 在 group(" + gconv.String(groupId) + ") 撤回了：\n"

		message = cqReplyRe.ReplaceAllString(message, "")
	}
	msg += message
	_, _ = service.Bot().SendMessage(ctx,
		"group", 0, notificationGroupId, msg, false)
	catch = true
	return
}
