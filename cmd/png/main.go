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
	arrowColor  = color.New(color.FgBlack).SprintFunc()

	okColor      = color.New(color.FgGreen).SprintFunc()
	timeoutColor = color.New(color.FgYellow).SprintFunc()
	errorColor   = color.New(color.FgHiRed).SprintFunc()

	elapsedColor = color.New(color.FgHiBlack).SprintFunc()
)

func main() {
	count := flag.IntP("count", "c", 0, "repeat count times (default: 0; means infinite repeat)")
	timeout := flag.DurationP("timeout", "t", 10*time.Second, "specify timeout")
	interval := flag.DurationP("interval", "i", 1*time.Second, "specify interval of ping iteration")
	noColor := flag.BoolP("no-color", "C", false, "disable color output")
	flag.Parse()

	color.NoColor = *noColor

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

	runner.Run()
}
