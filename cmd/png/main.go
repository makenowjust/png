package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/MakeNowJust/png"
	"github.com/fatih/color"
)

var (
	targetColor = color.New(color.FgWhite).Add(color.Bold).SprintfFunc()
	arrowColor  = color.New(color.FgBlack).SprintFunc()

	okColor      = color.New(color.FgGreen).SprintFunc()
	timeoutColor = color.New(color.FgYellow).SprintFunc()
	errorColor   = color.New(color.FgRed).Add(color.Bold).SprintFunc()
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

		pinger, err := png.Parse(target)
		if err != nil {
			log.Fatal(err)
		}

		pingers[i] = pinger
	}

	targetFmt := fmt.Sprintf("%%%ds", maxTargetLen)

	for i := uint(0); *count == 0 || i < *count; i++ {
		if i != 0 {
			<-time.After(*interval)
		}

		for i, pinger := range pingers {
			fmt.Printf("%s %s ", targetColor(targetFmt, targets[i]), arrowColor("->"))

			ctx, _ := context.WithTimeout(context.Background(), *timeout)

			done := make(chan error)
			go func() {
				done <- pinger.Ping(ctx)
			}()

			select {
			case <-ctx.Done():
				fmt.Println(timeoutColor("timeout"))
			case err := <-done:
				if err != nil {
					fmt.Println(errorColor("error"))
					fmt.Printf("  %v\n", err.Error())
				} else {
					fmt.Println(okColor("ok"))
				}
			}
		}
	}
}
