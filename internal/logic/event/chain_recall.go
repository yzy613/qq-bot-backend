package event

import (
	"context"
	"github.com/gogf/gf/v2/net/gtrace"
	"qq-bot-backend/internal/service"
)

func (s *sEvent) TryChainRecall(ctx context.Context) (catch bool) {
	msgIds, exist, err := service.Bot().GetCachedMessageContext(ctx,
		service.Bot().GetUserId(ctx),
		service.Bot().GetMsgId(ctx),
	)
	if err != nil || !exist {
		return
	}

	ctx, span := gtrace.NewSpan(ctx, "event.TryChainRecall")
	defer span.End()

	for _, msgId := range msgIds {
		service.Bot().RecallMessage(ctx, msgId)
	}

	catch = true
	return
}