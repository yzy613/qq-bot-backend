package command

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/util/gconv"
	"go.opentelemetry.io/otel/attribute"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
	"regexp"
	"strings"
)

type sCommand struct{}

func New() *sCommand {
	return &sCommand{}
}

func init() {
	service.RegisterCommand(New())
}

var (
	nextBranchRe      = regexp.MustCompile(`^(\S+)\s+([\s\S]+)$`)
	endBranchRe       = regexp.MustCompile(`^\S+$`)
	dualValueCmdEndRe = regexp.MustCompile(`^(\S+)\s+(\S+)$`)
)

func (s *sCommand) TryCommand(ctx context.Context, message string) (catch bool, retMsg string) {
	if !strings.HasPrefix(message, "/") {
		return
	}

	ctx, span := gtrace.NewSpan(ctx, "command.TryCommand")
	defer span.End()

	// 暂停状态时的权限校验
	userId := service.Bot().GetUserId(ctx)
	if !service.Process().IsBotProcessEnabled() &&
		!service.User().IsSystemTrustedUser(ctx, userId) {
		return
	}
	// 命令 log
	defer func() {
		if !catch {
			return
		}
		groupId := service.Bot().GetGroupId(ctx)
		span.SetAttributes(
			attribute.Int64("try_command.user_id", userId),
			attribute.Int64("try_command.group_id", groupId),
			attribute.String("try_command.command", message),
		)
		g.Log().Info(ctx,
			service.Bot().GetCardOrNickname(ctx)+"("+gconv.String(userId)+
				") in group("+gconv.String(groupId)+") send cmd "+message)
	}()
	cmd := strings.Replace(message, "/", "", 1)
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "list":
			// /list <>
			catch, retMsg = tryList(ctx, next[2])
		case "group":
			// /group <>
			catch, retMsg = tryGroup(ctx, next[2])
		case "namespace":
			// /namespace <>
			catch, retMsg = tryNamespace(ctx, next[2])
		case "user":
			// /user <>
			catch, retMsg = tryUser(ctx, next[2])
		case "raw":
			// 权限校验
			if !service.User().CanGetRawMsg(ctx, service.Bot().GetUserId(ctx)) {
				return
			}
			// span
			_, span := gtrace.NewSpan(ctx, "command.raw")
			// /raw <>
			catch, retMsg = true, next[2]
			span.End()
		case "broadcast":
			// /broadcast <>
			catch, retMsg = tryBroadcast(ctx, next[2])
		case "token":
			// /token <>
			catch, retMsg = tryToken(ctx, next[2])
		case "sys":
			// /sys <>
			catch, retMsg = trySys(ctx, next[2])
		case "model":
			// /model <>
			catch, retMsg = tryModelSet(ctx, next[2])
		}
	case endBranchRe.MatchString(cmd):
		// 权限校验
		if !service.User().IsSystemTrustedUser(ctx, service.Bot().GetUserId(ctx)) {
			return
		}
		switch cmd {
		case "status":
			// /status
			catch, retMsg = queryProcessStatus(ctx)
		case "version":
			// span
			_, span := gtrace.NewSpan(ctx, "command.version")
			// /version
			catch, retMsg = true, consts.Description
			span.End()
		case "continue":
			// /continue
			catch, retMsg = continueProcess(ctx)
		case "pause":
			// /pause
			catch, retMsg = pauseProcess(ctx)
		}
	}
	return
}
