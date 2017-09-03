package main

import (
	"time"

	"github.com/MakeNowJust/png"
)

type runner struct {
	count    int
	timeout  time.Duration
	interval time.Duration

	hookPingBefore func(target string)
	hookPingAfter  func(target, status string, elapsed time.Duration, err error)

	targets []string
	pingers []png.Pinger
}

func (r *runner) Run() {
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
		}
	}
}
