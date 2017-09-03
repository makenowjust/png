package main

import (
	"math"
	"time"

	"github.com/MakeNowJust/png"
)

type runner struct {
	count    int
	timeout  time.Duration
	interval time.Duration
	stats    string

	hookPingBefore  func(target string)
	hookPingAfter   func(target, status string, elapsed time.Duration, err error)
	hookStatsBefore func()
	hookStats       func(target string, ok, timeout, error, total int, min, max, average time.Duration)

	targets []string
	pingers []png.Pinger
}

func (r *runner) Run() {
	results := make([][]string, len(r.targets))
	durations := make([][]time.Duration, len(r.targets))

	for i := 0; r.count == 0 || i < r.count; i++ {
		if i != 0 {
			time.Sleep(r.interval)
		}

		for i, p := range r.pingers {
			if r.stats != "only" {
				r.hookPingBefore(r.targets[i])
			}

			elapsed, err := png.PingWithTimeout(p, r.timeout)
			var status string
			if err == nil {
				status = "ok"
			} else {
				if to, ok := err.(*png.Timeout); ok {
					status = "timeout"
					err = to.Err
				} else {
					status = "error"
				}
			}

			if r.stats != "only" {
				r.hookPingAfter(r.targets[i], status, elapsed, err)
			}
			results[i] = append(results[i], status)
			durations[i] = append(durations[i], elapsed)
		}
	}

	if r.stats == "all" {
		r.hookStatsBefore()
	}

	for i, target := range r.targets {
		total := len(results[i])
		if total == 0 {
			continue
		}

		ok := 0
		timeout := 0
		error := 0

		min := time.Duration(math.MaxInt64)
		max := time.Duration(0)
		average := time.Duration(0)

		for j, result := range results[i] {
			switch result {
			case "ok":
				ok += 1
			case "timeout":
				timeout += 1
			case "error":
				error += 1
			}

			elapsed := durations[i][j]

			if elapsed < min {
				min = elapsed
			}

			if elapsed > max {
				max = elapsed
			}

			average += elapsed / time.Duration(total)
		}

		if r.stats != "none" {
			r.hookStats(target, ok, timeout, error, total, min, max, average)
		}
	}
}
