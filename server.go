package main

import (
	"fmt"
	"net/http"

	"github.com/rodaine/x-files-rate/upstream"
)

func Serve(m RateLimiter, s *upstream.Service) {
	http.HandleFunc("/", m(*hertz, *burst, *wait, HelloWorld(s)))
	http.ListenAndServe("localhost:8080", nil)
}

func HelloWorld(s *upstream.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		switch err := s.Call(); err.(type) {
		case nil:
		// noop
		case upstream.QueueTimeout, upstream.WorkTimeout:
			rw.WriteHeader(http.StatusGatewayTimeout)
		default:
			rw.WriteHeader(http.StatusBadGateway)
		}

		fmt.Fprintln(rw, "Hello, World!")
	}
}
