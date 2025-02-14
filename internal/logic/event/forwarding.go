package event

import (
	"context"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/service"
	"sync"
)

func (s *sEvent) TryForward(ctx context.Context) (caught bool) {
	ctx, span := gtrace.NewSpan(ctx, "event.TryForward")
	defer span.End()

	if groupId := service.Bot().GetGroupId(ctx); groupId != 0 {
		if !service.Namespace().IsForwardingMatchGroupId(ctx, gconv.String(groupId)) {
			return
		}
	} else if userId := service.Bot().GetUserId(ctx); userId != 0 {
		if !service.Namespace().IsForwardingMatchUserId(ctx, gconv.String(userId)) {
			return
		}
	}
	caught = true

	wg := sync.WaitGroup{}
	defer wg.Wait()

	aliasList := service.Namespace().GetForwardingToAliasList(ctx)
	for alias := range aliasList {
		url, key := service.Namespace().GetForwardingTo(ctx, alias)
		if url == "" {
			return
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			service.Bot().Forward(ctx, url, key)
		}()
	}
	return
}
