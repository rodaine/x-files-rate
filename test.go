package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/tsenart/vegeta/lib"
)

func Test(rps uint64) {
	a := vegeta.NewAttacker()
	t := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    "http://127.0.0.1:8080/",
	})

	overall := &vegeta.Metrics{}
	success := &vegeta.Metrics{}
	failure := &vegeta.Metrics{}
	for r := range a.Attack(t, rps, 5*time.Second) {
		if r.Code == http.StatusOK {
			success.Add(r)
		} else {
			failure.Add(r)
		}
		overall.Add(r)
	}
	overall.Close()
	success.Close()
	failure.Close()

	fmt.Printf("%d, %.2f, %f, %f, %f, %f, %f, %f\n",
		rps, overall.Success,
		overall.Latencies.P50.Seconds()*1000, overall.Latencies.P95.Seconds()*1000,
		success.Latencies.P50.Seconds()*1000, success.Latencies.P95.Seconds()*1000,
		failure.Latencies.P50.Seconds()*1000, failure.Latencies.P95.Seconds()*1000)
}

type Latency struct {
}

type Results struct {
	Rate        float64
	SuccessRate float64
}
