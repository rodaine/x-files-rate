package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/config"
	"golang.org/x/time/rate"
)

type RateLimiter func(hertz, burst int, wait time.Duration, hf http.HandlerFunc) http.HandlerFunc

var targets = map[string]RateLimiter{
	"hello":     NoOp,
	"ticker":    Ticker,
	"fast":      Fast,
	"rate":      Rate,
	"tollbooth": Tollbooth,
}

func Resolve(target string) RateLimiter {
	m, ok := targets[target]
	if !ok {
		fmt.Println("unknown target:", target)
		os.Exit(1)
	}
	return m
}

func NoOp(hertz, burst int, wait time.Duration, hf http.HandlerFunc) http.HandlerFunc { return hf }

func Ticker(hertz, burst int, wait time.Duration, hf http.HandlerFunc) http.HandlerFunc {
	t, _ := tickerLimiter(hertz, burst)

	return func(w http.ResponseWriter, r *http.Request) {
		<-t
		hf(w, r)
	}
}

func Fast(hertz, burst int, wait time.Duration, hf http.HandlerFunc) http.HandlerFunc {
	t, _ := tickerLimiter(hertz, burst)

	return func(w http.ResponseWriter, r *http.Request) {

		select {
		case <-t: //noop
		case <-time.NewTimer(wait).C:
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		hf(w, r)
	}
}

func Rate(hertz, burst int, wait time.Duration, hf http.HandlerFunc) http.HandlerFunc {
	l := rate.NewLimiter(rate.Limit(hertz), burst)

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), wait)
		defer cancel()

		if err := l.Wait(ctx); err != nil {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		hf(w, r)
	}
}

func Tollbooth(hertz, burst int, wait time.Duration, hf http.HandlerFunc) http.HandlerFunc {
	l := config.NewLimiter(int64(burst), time.Second/time.Duration(hertz))
	return tollbooth.LimitFuncHandler(l, hf).ServeHTTP
}

func tickerLimiter(hertz, burst int) (<-chan time.Time, func()) {
	t := time.NewTicker(time.Second / time.Duration(hertz))
	c := make(chan time.Time, burst)

	for i := 0; i < burst; i++ {
		c <- time.Now()
	}

	go func() {
		for t := range t.C {
			select {
			case c <- t:
			default: // noop
			}
		}
	}()

	return c, t.Stop
}
