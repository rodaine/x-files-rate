package main

import (
	"flag"
	"syscall"

	"time"

	"github.com/rodaine/x-files-rate/upstream"
)

var (
	target = flag.String("target", "hello", "")
	rps    = flag.Uint64("rps", 10, "")

	concurrency = flag.Uint("concurrency", 10, "")
	mean        = flag.Duration("mean", 20*time.Millisecond, "")
	stddev      = flag.Duration("stddev", 5*time.Millisecond, "")
	timeout     = flag.Duration("timeout", 500*time.Millisecond, "")

	hertz = flag.Int("hertz", 425, "")
	burst = flag.Int("burst", 10, "")
	wait  = flag.Duration("wait", 75*time.Millisecond, "")
)

func main() {
	AdjustMaxFD()
	flag.Parse()

	m := Resolve(*target)
	s := upstream.New(*concurrency, *mean, *stddev, *timeout)

	go Serve(m, s)
	Test(*rps)
}

func AdjustMaxFD() {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	checkErr(err)

	rLimit.Max = 999999
	rLimit.Cur = 999999

	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	checkErr(err)
}
