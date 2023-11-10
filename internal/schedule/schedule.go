package schedule

import (
	"context"
	"reflect"
	"sync"
	"time"
)

var mutexMap = make(map[uintptr]*sync.Mutex)
var mapMutex sync.Mutex

func schedule(ctx context.Context, tick int, f func(context.Context)) {
	if tick < 2 {
		tick = 2
	}

	ticker := time.NewTicker(time.Duration(tick) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			m := mutex(f)
			m.Lock()

			go worker(m, time.Duration(tick-1)*time.Second, f)
		case <-ctx.Done():
			return
		}
	}
}

func worker(m *sync.Mutex, timeout time.Duration, f func(context.Context)) {
	defer m.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	go f(ctx)
	<-ctx.Done()
}

func mutex(f func(context.Context)) *sync.Mutex {
	fk := reflect.ValueOf(f).Pointer()

	mapMutex.Lock()
	defer mapMutex.Unlock()

	if _, exists := mutexMap[fk]; !exists {
		mutexMap[fk] = &sync.Mutex{}
	}
	return mutexMap[fk]
}
