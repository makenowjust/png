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

func main() {
	count := flag.IntP("count", "c", 0, "repeat count times (default: 0; means infinite repeat)")
	timeout := flag.DurationP("timeout", "t", 10*time.Second, "specify timeout")
	interval := flag.DurationP("interval", "i", 1*time.Second, "specify interval of ping iteration")
	noColor := flag.BoolP("no-color", "C", false, "disable color output")
	stats := flag.StringP("stats", "s", "", "decide to show statistics (default all; all/only/none)")
	format := flag.StringP("format", "f", "", "output format (default console; console/json)")

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

	if *format == "" {
		*format = "console"
	}

	if *format != "console" && *format != "json" {
		fmt.Fprintf(os.Stderr, "unknown format: %s\n", *format)
		flag.Usage()
		os.Exit(1)
	}

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	targets := flag.Args()
	pingers := make([]png.Pinger, len(targets))
	for i, target := range targets {
		pinger, err := png.Parse(target)
		if err != nil {
			log.Fatal(err)
		}
		pingers[i] = pinger
	}

	runner := &runner{
		count:    *count,
		timeout:  *timeout,
		interval: *interval,
		stats:    *stats,
		targets:  targets,
		pingers:  pingers,
	}

	switch *format {
	case "console":
		runner.HookConsole()
	case "json":
		runner.HookJSON()
	}

	runner.Run()
}
