package command

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
)

func tryGroupMessage(ctx context.Context, cmd string) (caught bool, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "enable":
			// /group message enable <>
			caught, retMsg = tryGroupMessageEnable(ctx, next[2])
		case "disable":
			// /group message disable <>
			caught, retMsg = tryGroupMessageDisable(ctx, next[2])
		case "set":
			// /group message set <>
			caught, retMsg = tryGroupMessageSet(ctx, next[2])
		case "rm":
			// /group message rm <>
			caught, retMsg = tryGroupMessageRemove(ctx, next[2])
		}
	}
	return
}

func tryGroupMessageEnable(ctx context.Context, cmd string) (caught bool, retMsg string) {
	switch {
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "anti-recall":
			// /group message enable anti-recall
			retMsg = service.Group().SetAntiRecallReturnRes(ctx, service.Bot().GetGroupId(ctx), true)
			caught = true
		}
	}
	return
}

func tryGroupMessageDisable(ctx context.Context, cmd string) (caught bool, retMsg string) {
	switch {
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "anti-recall":
			// /group message disable anti-recall
			retMsg = service.Group().SetAntiRecallReturnRes(ctx, service.Bot().GetGroupId(ctx), false)
			caught = true
		}
	}
	return
}

func tryGroupMessageSet(ctx context.Context, cmd string) (caught bool, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "notification":
			// /group message set notification <group_id>
			retMsg = service.Group().SetMessageNotificationReturnRes(ctx, service.Bot().GetGroupId(ctx), gconv.Int64(next[2]))
			caught = true
		}
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "only-anti-recall-member":
			// /group message set only-anti-recall-member
			retMsg = service.Group().SetOnlyAntiRecallMemberReturnRes(ctx, service.Bot().GetGroupId(ctx), true)
			caught = true
		}
	}
	return
}

func tryGroupMessageRemove(ctx context.Context, cmd string) (caught bool, retMsg string) {
	switch {
	case endBranchRe.MatchString(cmd):
		switch cmd {
		case "notification":
			// /group message rm notification
			retMsg = service.Group().RemoveMessageNotificationReturnRes(ctx, service.Bot().GetGroupId(ctx))
			caught = true
		case "only-anti-recall-member":
			// /group message rm only-anti-recall-member
			retMsg = service.Group().SetOnlyAntiRecallMemberReturnRes(ctx, service.Bot().GetGroupId(ctx), false)
			caught = true
		}
	}
	return
}
