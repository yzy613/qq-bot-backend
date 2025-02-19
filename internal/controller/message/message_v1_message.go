package message

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/util/gconv"
	"go.opentelemetry.io/otel/codes"
	"qq-bot-backend/internal/consts/errcode"
	"qq-bot-backend/internal/service"
	"qq-bot-backend/utility"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"

	"qq-bot-backend/api/message/v1"
)

func (c *ControllerV1) Message(ctx context.Context, req *v1.MessageReq) (res *v1.MessageRes, err error) {
	ctx, span := gtrace.NewSpan(ctx, "controller.Message")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	if req.Token == "" {
		// 忽视前置的 Bearer 或 Token 进行鉴权
		authorizations := strings.Fields(g.RequestFromCtx(ctx).Header.Get("Authorization"))
		if len(authorizations) < 2 {
			err = gerror.NewCode(errcode.Unauthorized)
			return
		}
		req.Token = authorizations[1]
	}
	// token 验证
	pass, tokenName, ownerId, botId := service.Token().IsCorrectToken(ctx, req.Token)
	if !pass {
		err = gerror.NewCode(errcode.Unauthorized)
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdmin(ctx, service.Namespace().GetGlobalNamespace(), ownerId) {
		if req.GroupId == 0 {
			err = gerror.NewCode(errcode.Forbidden)
			return
		}
		namespace := service.Group().GetNamespace(ctx, req.GroupId)
		if namespace == "" {
			err = gerror.NewCode(errcode.Forbidden)
			return
		}
		if !service.Namespace().IsNamespaceOwnerOrAdmin(ctx, namespace, ownerId) {
			err = gerror.NewCode(errcode.Forbidden)
			return
		}
	}
	// 记录访问时间
	service.Token().UpdateLoginTime(ctx, req.Token)
	// 加载 botId 对应的 botCtx
	botCtx := service.Bot().LoadConnectionPool(botId)
	if botCtx == nil {
		err = gerror.NewCode(errcode.BotNotConnected)
		return
	}
	// 规范请求参数
	if req.GroupId != 0 && req.UserId != 0 {
		req.UserId = 0
	}
	// for log
	{
		inner := struct {
			UserId  int64  `json:"user_id,omitempty"`
			GroupId int64  `json:"group_id,omitempty"`
			Message string `json:"message"`
		}{
			UserId:  req.UserId,
			GroupId: req.GroupId,
			Message: req.Message,
		}
		var innerStr string
		innerStr, err = sonic.MarshalString(inner)
		if err != nil {
			return
		}
		g.Log().Info(ctx, tokenName+" access successfully with "+innerStr)
	}
	// 限速 一分钟只能发送 7 条消息
	if limit, _ := utility.AutoLimit(ctx,
		"send_msg", gconv.String(req.UserId+req.GroupId), 7, time.Minute); limit {
		err = gerror.NewCode(errcode.TooMany)
		return
	}
	// send message
	if _, err = service.Bot().SendMessage(
		botCtx,
		service.Bot().GuessMsgType(req.GroupId),
		req.UserId,
		req.GroupId,
		req.Message, false,
	); err != nil {
		err = gerror.NewCode(errcode.InternalError, err.Error())
		return
	}
	return
}
