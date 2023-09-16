package controller

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"net/http"
	"qq-bot-backend/internal/service"
	"strings"
)

var (
	Bot = cBot{}
)

type cBot struct{}

func (c *cBot) Websocket(r *ghttp.Request) {
	ctx := r.Context()
	// 忽视前置的 Bearer 或 Token 进行鉴权
	authorizations := strings.Fields(r.Header.Get("Authorization"))
	if len(authorizations) < 2 {
		r.Response.WriteHeader(http.StatusForbidden)
		return
	}
	token := authorizations[1]
	var tokenName string
	if service.Cfg().IsEnabledDebug(ctx) {
		// token debug 验证模式
		var pass bool
		var name string
		pass, name, _ = service.Token().IsCorrectToken(ctx, token)
		// debug mode
		if !pass && token != service.Cfg().GetDebugToken(ctx) {
			r.Response.WriteHeader(http.StatusForbidden)
			return
		}
		if name == "" {
			tokenName = "debug"
		} else {
			tokenName = name
		}
	} else {
		// token 正常验证模式
		var pass bool
		pass, tokenName, _ = service.Token().IsCorrectToken(ctx, token)
		if !pass {
			r.Response.WriteHeader(http.StatusForbidden)
			return
		}
	}
	// 记录登录时间
	service.Token().UpdateLoginTime(ctx, token)
	// 升级 WebSocket 协议
	ws, err := r.WebSocket()
	if err != nil {
		return
	}
	g.Log().Info(ctx, tokenName+" Connected")
	// context 携带 WebSocket 对象
	ctx = service.Bot().CtxWithWebSocket(ctx, ws)
	// 并发 ws 写锁
	ctx = service.Bot().CtxNewWebSocketMutex(ctx)
	// 消息循环
	for {
		var wsReq []byte
		_, wsReq, err = ws.ReadMessage()
		if err != nil {
			g.Log().Info(ctx, tokenName+" Disconnected")
			return
		}
		// 异步处理 WebSocket 请求
		go service.Bot().Process(ctx, wsReq, service.Process().Process)
	}
}
