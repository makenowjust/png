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

	hookPingBefore  func(target string)
	hookPingAfter   func(target, status string, elapsed time.Duration, err error)
	hookStatsBefore func()
	hookStats       func(target string, n, succeed int, min, max, avg time.Duration)

	targets []string
	pingers []png.Pinger
}

func (r *runner) Run() {
	results := make([][]bool, len(r.targets))
	durations := make([][]time.Duration, len(r.targets))

	for i := 0; r.count == 0 || i < r.count; i++ {
		if i != 0 {
			time.Sleep(r.interval)
		}

		for i, p := range r.pingers {
			r.hookPingBefore(r.targets[i])

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

			r.hookPingAfter(r.targets[i], status, elapsed, err)
			results[i] = append(results[i], status == "ok")
			durations[i] = append(durations[i], elapsed)
		}
	}

	r.hookStatsBefore()

	for i, target := range r.targets {
		n := len(results[i])
		if n == 0 {
			continue
		}

		succeed := 0

		min := time.Duration(math.MaxInt64)
		max := time.Duration(0)
		avg := time.Duration(0)

		for j, result := range results[i] {
			if result {
				succeed += 1
			}

			elapsed := durations[i][j]

			if elapsed < min {
				min = elapsed
			}

			if elapsed > max {
				max = elapsed
			}

			avg += elapsed / time.Duration(n)
		}

		r.hookStats(target, n, succeed, min, max, avg)
	}
}
