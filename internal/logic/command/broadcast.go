package command

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
	"strings"
)

func tryBroadcast(ctx context.Context, cmd string) (catch bool, retMsg string) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "group":
			// /broadcast group <>
			if !nextBranchRe.MatchString(next[2]) {
				break
			}

			next := nextBranchRe.FindStringSubmatch(next[2])
			// /broadcast group <group_id> <>
			dstGroupId := gconv.Int64(next[1])
			userId := service.Bot().GetUserId(ctx)
			if dstNamespace := service.Group().GetNamespace(ctx, dstGroupId); dstNamespace == "" ||
				!service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, dstNamespace, userId) {
				break
			}

			suffix := "\n\nbroadcast from " + service.Bot().GetCardOrNickname(ctx) + "(" + gconv.String(userId) + ")"
			_ = service.Bot().SendMessage(ctx,
				service.Bot().GetMsgType(ctx),
				0,
				dstGroupId,
				strings.TrimSpace(next[2])+suffix,
				false,
			)
			catch = true
		}
	}
	return
}
