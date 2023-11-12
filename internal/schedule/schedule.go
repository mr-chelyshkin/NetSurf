package schedule

import (
	"context"
	"sync"
	"time"
)

func schedule(ctx context.Context, tick int, f func(context.Context)) {
	if tick < 2 {
		tick = 2
	}
	ticker := time.NewTicker(time.Duration(tick) * time.Second)
	defer ticker.Stop()

	var mutex sync.Mutex
	worker := func(timeout time.Duration, f func(context.Context)) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		mutex.Lock()
		go func() {
			defer mutex.Unlock()
			f(ctx)
		}()
		<-ctx.Done()
	}

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
