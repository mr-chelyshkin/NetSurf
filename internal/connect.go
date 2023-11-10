package internal

import (
	"context"
	"fmt"

	"github.com/mr-chelyshkin/NetSurf"
	"github.com/mr-chelyshkin/NetSurf/internal/schedule"
)

func connect(ctx context.Context, interrupt chan struct{}) {
	//output := make(chan string, 1)
	//defer close(output)

	networks := make(chan []map[string]string)
	defer close(networks)

	//ctx = context.WithValue(ctx, NetSurf.CtxKeyLoggerChannel, output)
	output, ok := ctx.Value(NetSurf.CtxKeyLoggerChannel).(chan string)
	if !ok {
		panic("AAAAAAA")
	}
	ctx, cancel := context.WithCancel(ctx)
	output <- "asd"
	go func() {
		for {
			select {
			case networks := <-networks:
				for _, network := range networks {
					output <- fmt.Sprintf("%s %s %s %s\n", network["ssid"], network["level"], network["quality"], network["freq"])
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	schedule.NetworkScan(ctx, networks)
	func() {
		select {
		case <-ctx.Done():
			return
		case <-interrupt:
			cancel()
			return
		}
	}()
}
