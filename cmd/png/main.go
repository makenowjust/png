package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/MakeNowJust/png"
)

func main() {
	count := flag.Uint("c", 0, "count")
	timeout := flag.Duration("t", 10*time.Second, "timeout")
	interval := flag.Duration("i", 1*time.Second, "interval")
	flag.Parse()

	targets := flag.Args()
	maxTargetLen := 0
	pingers := make([]png.Pinger, len(targets))
	for i, target := range targets {
		if maxTargetLen < len(target) {
			maxTargetLen = len(target)
		}

		pinger, err := png.NewPinger(target)
		if err != nil {
			log.Fatal(err)
		}

		pingers[i] = pinger
	}

	header := fmt.Sprintf("%%%ds -> ", maxTargetLen)

	for i := uint(0); *count == 0 || i < *count; i++ {
		if i != 0 {
			<-time.After(*interval)
		}

		for i, pinger := range pingers {
			fmt.Printf(header, targets[i])

			ctx, _ := context.WithTimeout(context.Background(), *timeout)

			done := make(chan error)
			go func() {
				done <- pinger.Ping(ctx)
			}()

			select {
			case <-ctx.Done():
				fmt.Println("timeout")
			case err := <-done:
				if err != nil {
					fmt.Println("error")
					fmt.Printf("  %v\n", err.Error())
				} else {
					fmt.Println("ok")
				}
			}
		}
	}
}
