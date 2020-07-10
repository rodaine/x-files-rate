# x-files-rate

An example program for the blog post [_The X-Files: Controlling Throughput with rate.Limiter_](http://rodaine.com/2017/05/x-files-time-rate-golang/), demonstrating various techniques for rate limiting an HTTP service via middleware.

## Use

### Bootstrap

```sh
go get github.com/rodaine/x-files-rate
```

### Run

Each rate-limiting middleware can be tested with ever-increasing RPS via [vegeta](https://github.com/tsenart/vegeta) using a customized vegeta test harness. The target names for the middleware can be found in [`middleware.go`](middleware.go). Results are printed to stdout and also stored in **./results/_{{target}}_.csv**.

```sh
./run.sh        # defaults to no rate-limiter
./run.sh ticker # targets the "ticker" middleware 
```

### Plot

The output CSV can be converted to SVG graphs using [gnuplot](http://www.gnuplot.info/). The SVG files are persisted in **./results/plots/_{{target}}_.svg**.

```sh
./plot.sh ticker
```
