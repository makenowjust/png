package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/MakeNowJust/png"
	"github.com/ttacon/chalk"
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

	header := fmt.Sprintf("%s%%%ds%s %s->%s ", chalk.White.NewStyle().WithBackground(chalk.ResetColor).WithTextStyle(chalk.Bold), maxTargetLen, chalk.Reset, chalk.Black, chalk.Reset)

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
				fmt.Println(chalk.Yellow.Color("timeout"))
			case err := <-done:
				if err != nil {
					fmt.Println(chalk.Red.Color("error"))
					fmt.Printf("  %v\n", err.Error())
				} else {
					fmt.Println(chalk.Green.Color("ok"))
				}
			}
		}
	}
}
