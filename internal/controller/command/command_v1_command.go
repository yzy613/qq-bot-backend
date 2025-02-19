package command

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"go.opentelemetry.io/otel/codes"
	"qq-bot-backend/internal/consts/errcode"
	"qq-bot-backend/internal/service"
	"qq-bot-backend/utility"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"

	"qq-bot-backend/api/command/v1"
)

func (c *ControllerV1) Command(ctx context.Context, req *v1.CommandReq) (res *v1.CommandRes, err error) {
	ctx, span := gtrace.NewSpan(ctx, "controller.Command")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	// 验证请求时间有效性
	{
		msgTime := gtime.New(time.Unix(req.Timestamp, 0))
		if diff := gtime.Now().Sub(msgTime); diff > 5*time.Second {
			err = gerror.NewCode(errcode.MessageExpired)
			return
		} else if diff < -5*time.Second {
			err = gerror.NewCode(errcode.TooEarly)
			return
		}
	}
	// 验证 token
	pass, tokenName, ownerId, botId := service.Token().IsCorrectToken(ctx, req.Token)
	if !pass {
		err = gerror.NewCode(errcode.Unauthorized)
		return
	}
	// 防止重放攻击
	if limit, _ := utility.AutoLimit(ctx,
		"api.command", req.Signature, 1, 10*time.Second); limit {
		err = gerror.NewCode(errcode.Conflict)
		return
	}
	// 验证签名
	{
		// 以 token+command+group_id+timestamp+message_sync+async 为原文，
		// 以 token_name 为 key 的 HmacSHA256 值的 base64 值
		s := req.Token + req.Command + gconv.String(req.GroupId) +
			gconv.String(req.Timestamp) + gconv.String(req.MessageSync) +
			gconv.String(req.Async)
		// HmacSHA256
		h := hmac.New(sha256.New, []byte(tokenName))
		h.Write([]byte(s))
		macBase64 := gbase64.Encode(h.Sum(nil))
		if !hmac.Equal(macBase64, []byte(req.Signature)) {
			err = gerror.NewCode(errcode.SignatureError)
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
	// 初始化内部请求
	innerReq := struct {
		ApiReq  struct{} `json:"api_req"`
		UserId  int64    `json:"user_id"`
		GroupId int64    `json:"group_id"`
	}{
		ApiReq:  struct{}{},
		UserId:  ownerId,
		GroupId: req.GroupId,
	}
	rawJson, err := sonic.MarshalString(innerReq)
	if err != nil {
		return
	}
	reqJson, _ := sonic.GetFromString(rawJson)
	botCtx = service.Bot().CtxWithReqJson(botCtx, &reqJson)
	g.Log().Info(ctx, tokenName+" access successfully with "+rawJson)
	var retMsg string
	// 异步执行
	if req.Async {
		go service.Command().TryCommand(botCtx, req.Command)
		retMsg = "async"
	} else {
		var caught bool
		caught, retMsg = service.Command().TryCommand(botCtx, req.Command)
		if !caught {
			err = gerror.NewCode(errcode.CommandNotFound)
			return
		}
	}
	// 响应
	res = &v1.CommandRes{
		Message: retMsg,
	}
	// 检查是否需要同步消息
	if !req.MessageSync || req.Async {
		return
	}
	if req.GroupId == 0 || !service.Group().IsBinding(botCtx, req.GroupId) {
		err = gerror.NewCode(errcode.GroupNotBound)
		return
	}
	if !service.Bot().IsGroupOwnerOrAdminOrSysTrusted(botCtx) {
		err = gerror.NewCode(errcode.Forbidden)
		return
	}
	// 限速 一分钟只能发送 5 条消息
	if limit, _ := utility.AutoLimit(ctx,
		"send_msg", gconv.String(req.GroupId), 5, time.Minute); limit {
		err = gerror.NewCode(errcode.TooMany)
		return
	}
	// 发送消息
	if _, err = service.Bot().SendMessage(
		botCtx,
		"",
		0,
		req.GroupId,
		retMsg,
		true,
	); err != nil {
		err = gerror.NewCode(errcode.InternalError, err.Error())
		return
	}
	return
}
