package schedule

import (
	"context"
	"time"
)

func schedule(ctx context.Context, tick int, f func(context.Context)) {
	if tick < 2 {
		tick = 2
	}
	ticker := time.NewTicker(time.Duration(tick) * time.Second)
	defer ticker.Stop()

	worker(time.Duration(tick-1)*time.Second, f)
	for {
		select {
		case <-ticker.C:
			worker(time.Duration(tick-1)*time.Second, f)
		case <-ctx.Done():
			return
		}
	}
}

func worker(timeout time.Duration, f func(context.Context)) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	go f(ctx)
	<-ctx.Done()
}
