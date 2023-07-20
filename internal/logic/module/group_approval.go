package module

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/consts"
	"qq-bot-backend/internal/service"
	"regexp"
)

func (s *sModule) TryApproveAddGroup(ctx context.Context) (catch bool) {
	// 获取当前 group approval 策略
	groupId := service.Bot().GetGroupId(ctx)
	process := service.Group().GetApprovalProcess(ctx, groupId)
	// 预处理
	if len(process) < 1 {
		// 没有入群审批策略，跳过审批功能
		return
	}
	comment := service.Bot().GetComment(ctx)
	// 默认通过审批
	pass := true
	// 局部变量
	userId := service.Bot().GetUserId(ctx)
	var extra string
	// 处理
	if _, ok := process[consts.McCmd]; ok {
		// mc 正版验证
		pass, extra = verifyMinecraftGenuine(ctx, comment)
	}
	if _, ok := process[consts.RegexpCmd]; ok && pass {
		// 正则表达式
		pass, extra = isMatchRegexp(ctx, groupId, comment)
	}
	if _, ok := process[consts.WhitelistCmd]; ok && pass {
		// 白名单
		pass = isInApprovalWhitelist(ctx, groupId, userId, extra)
	}
	if _, ok := process[consts.BlacklistCmd]; ok && pass {
		// 黑名单
		pass = isNotInApprovalBlacklist(ctx, groupId, userId, extra)
	}
	// 回执与日志
	var logMsg string
	if (!pass && service.Group().IsEnabledApprovalAutoReject(ctx, groupId)) ||
		(pass && service.Group().IsEnabledApprovalAutoPass(ctx, groupId)) {
		// 在不通过和启用自动通过的条件下发送审批回执
		// 审批请求回执
		service.Bot().ApproveJoinGroup(ctx,
			service.Bot().GetFlag(ctx),
			service.Bot().GetSubType(ctx),
			pass,
			"Auto-rejection")
		// 打印审批日志
		if pass {
			logMsg = fmt.Sprintf("approve user(%v) join group(%v) with %v",
				userId,
				groupId,
				comment)
		} else {
			logMsg = fmt.Sprintf("reject user(%v) join group(%v) with %v",
				userId,
				groupId,
				comment)
		}
	} else if pass {
		// 打印跳过同意日志
		logMsg = fmt.Sprintf("skip processing approve user(%v) join group(%v) with %v",
			userId,
			groupId,
			comment)
	} else if !pass {
		// 打印跳过拒绝日志
		logMsg = fmt.Sprintf("skip processing reject user(%v) join group(%v) with %v",
			userId,
			groupId,
			comment)
	}
	g.Log().Info(ctx, logMsg)
	// 通知
	notificationGroupId := service.Group().GetApprovalNotificationGroupId(ctx, groupId)
	if notificationGroupId > 0 {
		service.Bot().SendMessage(ctx,
			"group", 0, notificationGroupId, logMsg, true)
	}
	catch = true
	return
}

func isMatchRegexp(ctx context.Context, groupId int64, comment string) (match bool, matched string) {
	exp := service.Group().GetApprovalRegexp(ctx, groupId)
	// 匹配正则
	re, err := regexp.Compile(exp)
	if err != nil {
		g.Log().Warning(ctx, err)
		return
	}
	ans := re.FindStringSubmatch(comment)
	switch len(ans) {
	case 0:
	case 1:
		matched = ans[0]
		match = true
	default:
		// 读取第一个子表达式
		matched = ans[1]
		match = true
	}
	return
}

func verifyMinecraftGenuine(ctx context.Context, comment string) (genuine bool, uuid string) {
	// Minecraft 正版验证
	genuine, _, uuid, err := service.ThirdParty().QueryMinecraftGenuineUser(ctx, comment)
	if err != nil {
		g.Log().Notice(ctx, err)
	}
	return
}

func isInApprovalWhitelist(ctx context.Context, groupId, userId int64, extra string) (in bool) {
	// 获取白名单组
	whitelists := service.Group().GetApprovalWhitelists(ctx, groupId)
	for k := range whitelists {
		// 获取其中一个白名单
		whitelist := service.List().GetListData(ctx, k)
		if v, ok := whitelist[gconv.String(userId)]; ok {
			// userId 在白名单中
			if vv, okay := v.(string); okay {
				// 有额外验证信息
				if vv == extra {
					in = true
					return
				}
			} else {
				// 没有额外验证信息
				in = true
				return
			}
		}
		if extra == "" {
			// 没有额外验证信息则跳过反向验证
			continue
		}
		// 反向验证
		if v, ok := whitelist[extra]; ok {
			if vv, okay := v.(string); okay {
				if vv == gconv.String(userId) {
					in = true
					return
				}
			}
		}
	}
	return
}

func isNotInApprovalBlacklist(ctx context.Context, groupId, userId int64, extra string) (notIn bool) {
	// 默认不在黑名单内
	notIn = true
	// 获取黑名单组
	blacklists := service.Group().GetApprovalBlacklists(ctx, groupId)
	for k := range blacklists {
		// 获取其中一个黑名单
		blacklist := service.List().GetListData(ctx, k)
		if v, ok := blacklist[gconv.String(userId)]; ok {
			// userId 在黑名单中
			if vv, okay := v.(string); okay {
				// 有额外验证信息
				if vv == extra {
					notIn = false
					return
				}
			} else {
				// 没有额外验证信息
				notIn = false
				return
			}
		}
		if extra == "" {
			// 没有额外验证信息则跳过反向验证
			continue
		}
		// 反向验证
		if v, ok := blacklist[extra]; ok {
			if vv, okay := v.(string); okay {
				if vv == gconv.String(userId) {
					notIn = false
					return
				}
			}
		}
	}
	return
}