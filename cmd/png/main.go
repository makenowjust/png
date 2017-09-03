package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MakeNowJust/png"
	"github.com/fatih/color"
	flag "github.com/spf13/pflag"
)

var (
	targetColor = color.New(color.FgHiWhite).SprintfFunc()
	arrowColor  = color.New(color.FgBlack).SprintfFunc()

	okColor      = color.New(color.FgGreen).SprintfFunc()
	timeoutColor = color.New(color.FgYellow).SprintfFunc()
	errorColor   = color.New(color.FgHiRed).SprintfFunc()

	elapsedColor = color.New(color.FgHiBlack).SprintFunc()
)

func main() {
	count := flag.IntP("count", "c", 0, "repeat count times (default: 0; means infinite repeat)")
	timeout := flag.DurationP("timeout", "t", 10*time.Second, "specify timeout")
	interval := flag.DurationP("interval", "i", 1*time.Second, "specify interval of ping iteration")
	noColor := flag.BoolP("no-color", "C", false, "disable color output")
	stats := flag.StringP("stats", "s", "", "decide to show statistics (default all; all/only/none)")

	flag.Parse()
	color.NoColor = *noColor

	if *stats == "" {
		*stats = "all"
	}

	if *stats != "all" && *stats != "only" && *stats != "none" {
		fmt.Fprintf(os.Stderr, "unknown stats mode: %s\n", *stats)
		flag.Usage()
		os.Exit(1)
	}

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

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

	runner := &runner{
		count:    *count,
		timeout:  *timeout,
		interval: *interval,
		stats:    *stats,
		targets:  targets,
		pingers:  pingers,
	}

	runner.hookPingBefore = func(target string) {
		fmt.Printf("%s %s ", targetColor(targetFmt, target), arrowColor("->"))
	}

	runner.hookPingAfter = func(target, status string, elapsed time.Duration, err error) {
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

	runner.hookStatsBefore = func() {
		fmt.Println()
	}

	runner.hookStats = func(target string, ok, timeout, error, total int, min, max, avg time.Duration) {
		color := okColor
		if ok != total {
			if timeout > error {
				color = timeoutColor
			} else {
				color = errorColor
			}
		}

		fmt.Printf("%s: %s, min/max/avg = %12s/%12s/%12s\n",
			targetColor(targetFmt, target),
			color("ok/timeout/error/total = %2d/%2d/%2d/%2d", ok, timeout, error, total),
			min, max, avg)
	}

	runner.Run()
}
