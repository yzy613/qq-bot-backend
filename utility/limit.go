package utility

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"time"
)

func AutoLimit(ctx context.Context,
	kind, key string,
	limitTimes int,
	duration time.Duration,
) (limited bool, times int) {
	// 缓存键名
	cacheKey := "LimitTimes_" + kind + "_" + key

	// Try to get the cache value, or set it with the default value if it doesn't exist
	timesVar, err := gcache.GetOrSet(ctx, cacheKey, 1, duration)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}

	// Convert to int
	times = timesVar.Int()

	// Update the cache value by incrementing it
	if _, _, err = gcache.Update(ctx, cacheKey, times+1); err != nil {
		g.Log().Error(ctx, err)
		return
	}

	// Check if the times exceed the limit
	if times > limitTimes {
		limited = true
	}
	return
}
