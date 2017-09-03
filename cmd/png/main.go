package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/MakeNowJust/png"
	"github.com/fatih/color"
)

var (
	targetColor = color.New(color.FgHiWhite).SprintfFunc()
	arrowColor  = color.New(color.FgBlack).SprintFunc()

	okColor      = color.New(color.FgGreen).SprintFunc()
	timeoutColor = color.New(color.FgYellow).SprintFunc()
	errorColor   = color.New(color.FgHiRed).SprintFunc()

	elapsedColor = color.New(color.FgHiBlack).SprintFunc()
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
			time.Sleep(*interval)
		}

		for i, p := range pingers {
			fmt.Printf("%s %s ", targetColor(targetFmt, targets[i]), arrowColor("->"))

			elapsed, err := png.PingWithTimeout(p, *timeout)
			if err == nil {
				fmt.Printf("%s %s\n", okColor("ok     "), elapsedColor(elapsed))
			} else {
				switch err.(type) {
				case *png.Timeout:
					fmt.Printf("%s %s\n", timeoutColor("timeout"), elapsedColor(elapsed))
				default:
					fmt.Printf("%s %s\n  %v\n", errorColor("error  "), elapsedColor(elapsed), err)
				}
			}
		}
	}
}
