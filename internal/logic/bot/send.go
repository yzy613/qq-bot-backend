package bot

import (
	"context"
	"errors"
	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gorilla/websocket"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"qq-bot-backend/utility/segment"
	"sync"
)

func (s *sBot) SendMessage(ctx context.Context,
	messageType string,
	userId, groupId int64,
	msg string,
	plain bool,
) (messageId int64, err error) {
	// 参数校验
	if userId == 0 && groupId == 0 {
		return 0, errors.New("userId 和 groupId 不能同时为 0")
	}

	ctx, span := gtrace.NewSpan(ctx, "bot.SendMessage")
	defer span.End()
	span.SetAttributes(attribute.String("send_message.message", msg))
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	if groupId != 0 {
		userId = 0
		span.SetAttributes(attribute.Int64("send_message.group_id", groupId))
	} else {
		span.SetAttributes(attribute.Int64("send_message.user_id", userId))
	}

	// echo sign
	echoSign := s.generateEchoSignWithTrace(ctx)
	// 参数
	req := struct {
		Action string `json:"action"`
		Echo   string `json:"echo"`
		Params struct {
			MessageType string `json:"message_type,omitempty"`
			Message     any    `json:"message"`
			AutoEscape  bool   `json:"auto_escape,omitempty"`
			UserId      int64  `json:"user_id,omitempty"`
			GroupId     int64  `json:"group_id,omitempty"`
		} `json:"params"`
	}{
		Action: "send_msg",
		Echo:   echoSign,
		Params: struct {
			MessageType string `json:"message_type,omitempty"`
			Message     any    `json:"message"`
			AutoEscape  bool   `json:"auto_escape,omitempty"`
			UserId      int64  `json:"user_id,omitempty"`
			GroupId     int64  `json:"group_id,omitempty"`
		}{
			MessageType: messageType,
			Message:     msg,
			AutoEscape:  plain,
			UserId:      userId,
			GroupId:     groupId,
		},
	}
	// message segment
	if s.isMessageSegment(ctx) {
		if plain {
			req.Params.Message = segment.NewTextSegments(msg)
			req.Params.AutoEscape = false
		} else {
			req.Params.Message = segment.ParseMessage(msg)
		}
	}

	reqJson, err := sonic.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// callback
	wg := sync.WaitGroup{}
	defer wg.Wait()
	wgDone := sync.OnceFunc(wg.Done)
	wg.Add(1)
	callback := func(ctx context.Context, asyncCtx context.Context) {
		defer wgDone()
		if err = s.defaultEchoHandler(asyncCtx); err != nil {
			return
		}
		messageId = s.getMessageIdFromData(asyncCtx)
	}
	timeout := func(ctx context.Context) {
		defer wgDone()
		err = errors.New("echo timeout")
	}
	// echo
	if err = s.pushEchoCache(ctx, echoSign, callback, timeout); err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 发送响应
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJson); err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	return
}

// SendPlainMsg 适用于*事件*触发的消息发送
func (s *sBot) SendPlainMsg(ctx context.Context, msg string) {
	_, _ = s.SendMessage(ctx, s.GetMsgType(ctx), s.GetUserId(ctx), s.GetGroupId(ctx), msg, true)
}

// SendMsg 适用于*事件*触发的消息发送
func (s *sBot) SendMsg(ctx context.Context, msg string) {
	_, _ = s.SendMessage(ctx, s.GetMsgType(ctx), s.GetUserId(ctx), s.GetGroupId(ctx), msg, false)
}

func (s *sBot) SendMsgIfNotApiReq(ctx context.Context, msg string, richText ...bool) {
	if s.isApiReq(ctx) {
		return
	}
	s.SendMsgCacheContext(ctx, msg, richText...)
}

// SendMsgCacheContext 适用于*非事件*触发的消息发送
func (s *sBot) SendMsgCacheContext(ctx context.Context, msg string, richText ...bool) {
	plain := true
	if len(richText) > 0 && richText[0] {
		plain = false
	}
	sentMsgId, err := s.SendMessage(ctx, s.GetMsgType(ctx), s.GetUserId(ctx), s.GetGroupId(ctx), msg, plain)
	if err != nil {
		return
	}
	_ = s.CacheMessageContext(ctx, sentMsgId)
}

func (s *sBot) SendFileToGroup(ctx context.Context, groupId int64, filePath, name, folder string) {
	ctx, span := gtrace.NewSpan(ctx, "bot.SendFileToGroup")
	defer span.End()
	span.SetAttributes(
		attribute.Int64("send_file_to_group.group_id", groupId),
		attribute.String("send_file_to_group.file_path", filePath),
		attribute.String("send_file_to_group.name", name),
		attribute.String("send_file_to_group.folder", folder),
	)
	var err error
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	// echo sign
	echoSign := s.generateEchoSignWithTrace(ctx)
	// 参数
	req := struct {
		Action string `json:"action"`
		Echo   string `json:"echo"`
		Params struct {
			GroupId int64  `json:"group_id"`
			File    string `json:"file"`
			Name    string `json:"name"`
			Folder  string `json:"folder,omitempty"`
		} `json:"params"`
	}{
		Action: "upload_group_file",
		Echo:   echoSign,
		Params: struct {
			GroupId int64  `json:"group_id"`
			File    string `json:"file"`
			Name    string `json:"name"`
			Folder  string `json:"folder,omitempty"`
		}{
			GroupId: groupId,
			File:    filePath,
			Name:    name,
			Folder:  folder,
		},
	}
	reqJson, err := sonic.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// callback
	wg := sync.WaitGroup{}
	defer wg.Wait()
	wgDone := sync.OnceFunc(wg.Done)
	wg.Add(1)
	callback := func(ctx context.Context, asyncCtx context.Context) {
		defer wgDone()
		if err = s.defaultEchoHandler(asyncCtx); err != nil {
			s.SendMsgIfNotApiReq(ctx, err.Error())
			return
		}
	}
	timeout := func(ctx context.Context) {
		defer wgDone()
		err = errors.New("echo timeout")
		s.SendMsgIfNotApiReq(ctx, "上传至群文件超时")
	}
	// echo
	if err = s.pushEchoCache(ctx, echoSign, callback, timeout); err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJson); err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) SendFileToUser(ctx context.Context, userId int64, filePath, name string) {
	ctx, span := gtrace.NewSpan(ctx, "bot.SendFileToUser")
	defer span.End()
	span.SetAttributes(
		attribute.Int64("send_file_to_user.user_id", userId),
		attribute.String("send_file_to_user.file_path", filePath),
		attribute.String("send_file_to_user.name", name),
	)
	var err error
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	// echo sign
	echoSign := s.generateEchoSignWithTrace(ctx)
	// 参数
	req := struct {
		Action string `json:"action"`
		Echo   string `json:"echo"`
		Params struct {
			UserId int64  `json:"user_id"`
			File   string `json:"file"`
			Name   string `json:"name"`
		} `json:"params"`
	}{
		Action: "upload_private_file",
		Echo:   echoSign,
		Params: struct {
			UserId int64  `json:"user_id"`
			File   string `json:"file"`
			Name   string `json:"name"`
		}{
			UserId: userId,
			File:   filePath,
			Name:   name,
		},
	}
	reqJson, err := sonic.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// callback
	wg := sync.WaitGroup{}
	defer wg.Wait()
	wgDone := sync.OnceFunc(wg.Done)
	wg.Add(1)
	callback := func(ctx context.Context, asyncCtx context.Context) {
		defer wgDone()
		if err = s.defaultEchoHandler(asyncCtx); err != nil {
			s.SendMsgIfNotApiReq(ctx, err.Error())
			return
		}
	}
	timeout := func(ctx context.Context) {
		defer wgDone()
		err = errors.New("echo timeout")
		s.SendMsgIfNotApiReq(ctx, "上传文件至私聊超时")
	}
	// echo
	if err = s.pushEchoCache(ctx, echoSign, callback, timeout); err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJson); err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) SendFile(ctx context.Context, filePath, name string) {
	if groupId := s.GetGroupId(ctx); groupId != 0 {
		s.SendFileToGroup(ctx, groupId, filePath, name, "")
		return
	}
	s.SendFileToUser(ctx, s.GetUserId(ctx), filePath, name)
}

func (s *sBot) UploadFile(ctx context.Context, url string) (filePath string, err error) {
	ctx, span := gtrace.NewSpan(ctx, "bot.UploadFile")
	defer span.End()
	span.SetAttributes(attribute.String("upload_file.url", url))
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	// echo sign
	echoSign := s.generateEchoSignWithTrace(ctx)
	// 参数
	req := struct {
		Action string `json:"action"`
		Echo   string `json:"echo"`
		Params struct {
			URL string `json:"url"`
		} `json:"params"`
	}{
		Action: "download_file",
		Echo:   echoSign,
		Params: struct {
			URL string `json:"url"`
		}{
			URL: url,
		},
	}
	reqJson, err := sonic.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// callback
	wg := sync.WaitGroup{}
	defer wg.Wait()
	wgDone := sync.OnceFunc(wg.Done)
	wg.Add(1)
	callback := func(ctx context.Context, asyncCtx context.Context) {
		defer wgDone()
		if err = s.defaultEchoHandler(asyncCtx); err != nil {
			s.SendMsgIfNotApiReq(ctx, err.Error())
			return
		}
		filePath = s.getFileFromData(asyncCtx)
	}
	timeout := func(ctx context.Context) {
		defer wgDone()
		err = errors.New("echo timeout")
	}
	// echo
	if err = s.pushEchoCache(ctx, echoSign, callback, timeout); err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJson); err != nil {
		g.Log().Warning(ctx, err)
	}
	return
}

func (s *sBot) ApproveJoinGroup(ctx context.Context, flag, subType string, approve bool, reason string) {
	ctx, span := gtrace.NewSpan(ctx, "bot.ApproveJoinGroup")
	defer span.End()
	span.SetAttributes(
		attribute.String("approve_join_group.flag", flag),
		attribute.String("approve_join_group.sub_type", subType),
		attribute.Bool("approve_join_group.approve", approve),
		attribute.String("approve_join_group.reason", reason),
	)
	var err error
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	// 参数校验
	if approve {
		reason = ""
	}
	// echo sign
	echoSign := s.generateEchoSignWithTrace(ctx)
	// 参数
	req := struct {
		Action string `json:"action"`
		Echo   string `json:"echo"`
		Params struct {
			Flag    string `json:"flag"`
			SubType string `json:"sub_type"`
			Approve bool   `json:"approve"`
			Reason  string `json:"reason,omitempty"`
		} `json:"params"`
	}{
		Action: "set_group_add_request",
		Echo:   echoSign,
		Params: struct {
			Flag    string `json:"flag"`
			SubType string `json:"sub_type"`
			Approve bool   `json:"approve"`
			Reason  string `json:"reason,omitempty"`
		}{
			Flag:    flag,
			SubType: subType,
			Approve: approve,
			Reason:  reason,
		},
	}
	reqJson, err := sonic.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// callback
	wg := sync.WaitGroup{}
	defer wg.Wait()
	wgDone := sync.OnceFunc(wg.Done)
	wg.Add(1)
	callback := func(ctx context.Context, asyncCtx context.Context) {
		defer wgDone()
		if err = s.defaultEchoHandler(asyncCtx); err != nil {
			return
		}
	}
	timeout := func(ctx context.Context) {
		defer wgDone()
		err = errors.New("echo timeout")
	}
	// echo
	if err = s.pushEchoCache(ctx, echoSign, callback, timeout); err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 发送响应
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJson); err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) SetModel(ctx context.Context, model string) {
	ctx, span := gtrace.NewSpan(ctx, "bot.SetModel")
	defer span.End()
	span.SetAttributes(attribute.String("set_model.model", model))
	var err error
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	// echo sign
	echoSign := s.generateEchoSignWithTrace(ctx)
	// 参数
	req := struct {
		Action string `json:"action"`
		Echo   string `json:"echo"`
		Params struct {
			Model     string `json:"model"`
			ModelShow string `json:"model_show"`
		} `json:"params"`
	}{
		Action: "_set_model_show",
		Echo:   echoSign,
		Params: struct {
			Model     string `json:"model"`
			ModelShow string `json:"model_show"`
		}{
			Model:     model,
			ModelShow: model,
		},
	}
	reqJson, err := sonic.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// callback
	wg := sync.WaitGroup{}
	defer wg.Wait()
	wgDone := sync.OnceFunc(wg.Done)
	wg.Add(1)
	callback := func(ctx context.Context, asyncCtx context.Context) {
		defer wgDone()
		if err = s.defaultEchoHandler(asyncCtx); err != nil {
			s.SendMsgIfNotApiReq(ctx, err.Error())
			return
		}
		s.SendMsgIfNotApiReq(ctx, "已更改机型为 '"+model+"'")
	}
	timeout := func(ctx context.Context) {
		defer wgDone()
		err = errors.New("echo timeout")
		s.SendMsgIfNotApiReq(ctx, "更改机型超时")
	}
	// echo
	if err = s.pushEchoCache(ctx, echoSign, callback, timeout); err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 发送响应
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJson); err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) RecallMessage(ctx context.Context, messageId int64) {
	ctx, span := gtrace.NewSpan(ctx, "bot.RecallMessage")
	defer span.End()
	span.SetAttributes(attribute.Int64("recall_message.message_id", messageId))
	var err error
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	// echo sign
	echoSign := s.generateEchoSignWithTrace(ctx)
	// 参数
	req := struct {
		Action string `json:"action"`
		Echo   string `json:"echo"`
		Params struct {
			MessageId int64 `json:"message_id"`
		} `json:"params"`
	}{
		Action: "delete_msg",
		Echo:   echoSign,
		Params: struct {
			MessageId int64 `json:"message_id"`
		}{
			MessageId: messageId,
		},
	}
	reqJson, err := sonic.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// callback
	wg := sync.WaitGroup{}
	defer wg.Wait()
	wgDone := sync.OnceFunc(wg.Done)
	wg.Add(1)
	callback := func(ctx context.Context, asyncCtx context.Context) {
		defer wgDone()
		if err = s.defaultEchoHandler(asyncCtx); err != nil {
			return
		}
	}
	timeout := func(ctx context.Context) {
		defer wgDone()
		err = errors.New("echo timeout")
	}
	// echo
	if err = s.pushEchoCache(ctx, echoSign, callback, timeout); err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJson); err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) MutePrototype(ctx context.Context, groupId, userId int64, seconds int) {
	ctx, span := gtrace.NewSpan(ctx, "bot.MutePrototype")
	defer span.End()
	span.SetAttributes(
		attribute.Int64("mute_prototype.group_id", groupId),
		attribute.Int64("mute_prototype.user_id", userId),
		attribute.Int("mute_prototype.seconds", seconds),
	)
	var err error
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	// 参数校验
	if seconds > 2591940 {
		// 不大于 29 天 23 小时 59 分钟
		// (30*24*60-1)*60=2591940 秒
		seconds = 2591940
	}
	// echo sign
	echoSign := s.generateEchoSignWithTrace(ctx)
	// 参数
	req := struct {
		Action string `json:"action"`
		Echo   string `json:"echo"`
		Params struct {
			GroupId  int64 `json:"group_id"`
			UserId   int64 `json:"user_id"`
			Duration int   `json:"duration"`
		} `json:"params"`
	}{
		Action: "set_group_ban",
		Echo:   echoSign,
		Params: struct {
			GroupId  int64 `json:"group_id"`
			UserId   int64 `json:"user_id"`
			Duration int   `json:"duration"`
		}{
			GroupId:  groupId,
			UserId:   userId,
			Duration: seconds,
		},
	}
	reqJson, err := sonic.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// callback
	wg := sync.WaitGroup{}
	defer wg.Wait()
	wgDone := sync.OnceFunc(wg.Done)
	wg.Add(1)
	callback := func(ctx context.Context, asyncCtx context.Context) {
		defer wgDone()
		if err = s.defaultEchoHandler(asyncCtx); err != nil {
			return
		}
	}
	timeout := func(ctx context.Context) {
		defer wgDone()
		err = errors.New("echo timeout")
	}
	// echo
	if err = s.pushEchoCache(ctx, echoSign, callback, timeout); err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJson); err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) Mute(ctx context.Context, seconds int) {
	s.MutePrototype(ctx, s.GetGroupId(ctx), s.GetUserId(ctx), seconds)
}

func (s *sBot) SetGroupCard(ctx context.Context, groupId, userId int64, card string) {
	ctx, span := gtrace.NewSpan(ctx, "bot.SetGroupCard")
	defer span.End()
	span.SetAttributes(
		attribute.Int64("set_group_card.group_id", groupId),
		attribute.Int64("set_group_card.user_id", userId),
		attribute.String("set_group_card.card", card),
	)
	var err error
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	// echo sign
	echoSign := s.generateEchoSignWithTrace(ctx)
	// 参数
	req := struct {
		Action string `json:"action"`
		Echo   string `json:"echo"`
		Params struct {
			GroupId int64  `json:"group_id"`
			UserId  int64  `json:"user_id"`
			Card    string `json:"card"`
		} `json:"params"`
	}{
		Action: "set_group_card",
		Echo:   echoSign,
		Params: struct {
			GroupId int64  `json:"group_id"`
			UserId  int64  `json:"user_id"`
			Card    string `json:"card"`
		}{
			GroupId: groupId,
			UserId:  userId,
			Card:    card,
		},
	}
	reqJson, err := sonic.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// callback
	wg := sync.WaitGroup{}
	defer wg.Wait()
	wgDone := sync.OnceFunc(wg.Done)
	wg.Add(1)
	callback := func(ctx context.Context, asyncCtx context.Context) {
		defer wgDone()
		if err = s.defaultEchoHandler(asyncCtx); err != nil {
			return
		}
	}
	timeout := func(ctx context.Context) {
		defer wgDone()
		err = errors.New("echo timeout")
	}
	// echo
	if err = s.pushEchoCache(ctx, echoSign, callback, timeout); err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJson); err != nil {
		g.Log().Warning(ctx, err)
	}
}

func (s *sBot) Kick(ctx context.Context, groupId, userId int64, reject ...bool) {
	ctx, span := gtrace.NewSpan(ctx, "bot.Kick")
	defer span.End()
	span.SetAttributes(
		attribute.Int64("kick.group_id", groupId),
		attribute.Int64("kick.user_id", userId),
	)
	var err error
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	// echo sign
	echoSign := s.generateEchoSignWithTrace(ctx)
	// 参数
	req := struct {
		Action string `json:"action"`
		Echo   string `json:"echo"`
		Params struct {
			GroupId          int64 `json:"group_id"`
			UserId           int64 `json:"user_id"`
			RejectAddRequest bool  `json:"reject_add_request,omitempty"`
		} `json:"params"`
	}{
		Action: "set_group_kick",
		Echo:   echoSign,
		Params: struct {
			GroupId          int64 `json:"group_id"`
			UserId           int64 `json:"user_id"`
			RejectAddRequest bool  `json:"reject_add_request,omitempty"`
		}{
			GroupId:          groupId,
			UserId:           userId,
			RejectAddRequest: false,
		},
	}
	if len(reject) > 0 && reject[0] {
		req.Params.RejectAddRequest = true
	}
	reqJson, err := sonic.Marshal(req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// callback
	wg := sync.WaitGroup{}
	defer wg.Wait()
	wgDone := sync.OnceFunc(wg.Done)
	wg.Add(1)
	callback := func(ctx context.Context, asyncCtx context.Context) {
		defer wgDone()
		if err = s.defaultEchoHandler(asyncCtx); err != nil {
			return
		}
	}
	timeout := func(ctx context.Context) {
		defer wgDone()
		err = errors.New("echo timeout")
	}
	// echo
	if err = s.pushEchoCache(ctx, echoSign, callback, timeout); err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if err = s.writeMessage(ctx, websocket.TextMessage, reqJson); err != nil {
		g.Log().Warning(ctx, err)
	}
}
