package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type result struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type ping struct {
	Target  string        `json:"target"`
	Status  string        `json:"status"`
	Elapsed time.Duration `json:"elapsed"`
	Err     string        `json:"err,omitempty"`
}

type stats struct {
	Target  string        `json:"target"`
	Ok      int           `json:"ok"`
	Timeout int           `json:"timeout"`
	Error   int           `json:"error"`
	Total   int           `json:"total"`
	Min     time.Duration `json:"min"`
	Max     time.Duration `json:"max"`
	Average time.Duration `json:"average"`
}

func (r *runner) HookJSON() {
	errString := func(err error) string {
		if err == nil {
			return ""
		}
		return err.Error()
	}

	r.hookPingBefore = func(target string) {}
	r.hookStatsBefore = func() {}

	r.hookPingAfter = func(target, status string, elapsed time.Duration, err error) {
		json, err := json.Marshal(&result{
			Type: "ping",
			Payload: &ping{
				Target:  target,
				Status:  status,
				Elapsed: elapsed,
				Err:     errString(err),
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		os.Stdout.Write(json)
		fmt.Println()
	}

	r.hookStats = func(target string, ok, timeout, error, total int, min, max, average time.Duration) {
		json, err := json.Marshal(&result{
			Type: "stats",
			Payload: &stats{
				Target:  target,
				Ok:      ok,
				Timeout: timeout,
				Error:   error,
				Total:   total,
				Min:     min,
				Max:     max,
				Average: average,
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		os.Stdout.Write(json)
		fmt.Println()
	}
}
