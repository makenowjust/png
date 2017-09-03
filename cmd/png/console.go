package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

var (
	targetColor = color.New(color.FgHiWhite).SprintfFunc()
	arrowColor  = color.New(color.FgBlack).SprintfFunc()

	okColor      = color.New(color.FgGreen).SprintfFunc()
	timeoutColor = color.New(color.FgYellow).SprintfFunc()
	errorColor   = color.New(color.FgHiRed).SprintfFunc()

	elapsedColor = color.New(color.FgHiBlack).SprintFunc()
)

func (r *runner) HookConsole() {
	maxTargetLen := 0
	for _, target := range r.targets {
		if maxTargetLen < len(target) {
			maxTargetLen = len(target)
		}
	}

	targetFmt := fmt.Sprintf("%%%ds", maxTargetLen)

	r.hookPingBefore = func(target string) {
		fmt.Printf("%s %s ", targetColor(targetFmt, target), arrowColor("->"))
	}

	r.hookPingAfter = func(target, status string, elapsed time.Duration, err error) {
		padStatus := fmt.Sprintf("%-7s", status)

		switch status {
		case "ok":
			fmt.Printf("%s %s\n", okColor(padStatus), elapsedColor(elapsed))
		case "timeout":
			fmt.Printf("%s %s\n", timeoutColor(padStatus), elapsedColor(elapsed))
		case "error":
			fmt.Printf("%s %s\n  %v\n", errorColor(padStatus), elapsedColor(elapsed), err)
		}
	}

	r.hookStatsBefore = func() {
		fmt.Println()
	}

	r.hookStats = func(target string, ok, timeout, error, total int, min, max, average time.Duration) {
		color := okColor
		if ok != total {
			if timeout > error {
				color = timeoutColor
			} else {
				color = errorColor
			}
		}

		fmt.Printf("%s: %s, min/max/average = %12s/%12s/%12s\n",
			targetColor(targetFmt, target),
			color("ok/timeout/error/total = %2d/%2d/%2d/%2d", ok, timeout, error, total),
			min, max, average)
	}
}
