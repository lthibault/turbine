# turbine

[![Godoc Reference](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/lthibault/turbine/pkg) [![Go Report Card](https://goreportcard.com/badge/github.com/lthibault/turbine?style=flat-square)](https://goreportcard.com/report/github.com/lthibault/turbine)

Lock-free pipeline for inter-thread message passing.

## Overview

Turbine is a concurrency data structure that allows for lock-free message passing between threads.

It's a special kind of ring-buffer that supports batched reading and writng, as well as multi-stage, concurrent pipelines.  It is allocation and garbage free, and exhibits latency on the order of nanoseconds.

Turbine is based heavily on [go-disruptor](https://github.com/smartystreets/go-disruptor). It features a simplified architecture, and fixes the data-races inherent in go-disruptor's
design.

## Why Turbine

Turbine is not suitable for all scenarios.  **When in doubt, use channels.**  With Turbine, you lose many of Go's most useful concurrency features:

1. `select` statements
2. synchronous message passing (i.e.: async only)
3. synchronous close semantics

There are two situations in which Turbine is useful.

The first is when the concurrency problem can be modeled as a concurrent pipeline with independent stages that run in strict order.  This kind of problem is (in our opinion) more naturally expressed with `Turbine` than with channels, though your mileage may vary.

The second use case is in ultra low-latency environments where jitter must be kept to a minimum.

## Benchmarks

Benchmarks can be run as follows:

```bash
$ go test -bench=. ./...

goos: darwin
goarch: amd64
pkg: github.com/lthibault/turbine/pkg
BenchmarkTurbine/ReserveOne-4           50000000                37.2 ns/op             0 B/op          0 allocs/op
BenchmarkTurbine/ReserveMany-4          300000000                6.33 ns/op            0 B/op          0 allocs/op
BenchmarkChan/Unbuffered-4               5000000               361 ns/op               0 B/op          0 allocs/op
BenchmarkChan/Buffered-1-4               5000000               261 ns/op               0 B/op          0 allocs/op
BenchmarkChan/Buffered-16-4             10000000               113 ns/op               0 B/op          0 allocs/op
PASS
ok      github.com/lthibault/turbine/pkg        9.448s
```

Suggestions for improving benchmark accuracy are welcome.

## Quickstart

The following demonstrates the use of `turbine`.  Note the absence of `interface{}`!

```go
package main

import (
    "log"

    turbine "github.com/lthibault/turbine/pkg"
)

const (
    cap  = 1024 // must be a power of 2
    mask = cap - 1
)

var ring = [cap]packet{}

type packet struct {
    Value int
}

type multiplier struct{}

func (multiplier) Consume(lower, upper int64) {
    for seq := lower; seq <= upper; seq++ {
        ring[seq&mask].Value *= 100
    }
}

type logger struct{}

func (logger) Consume(lower, upper int64) {
    var msg packet
    for seq := lower; seq <= upper; seq++ {
        msg = ring[seq&mask]
        log.Println(msg.Value)
    }
}

func main() {
    t := turbine.New(cap, multiplier{}, logger{})
    t.Start()
    defer t.Stop()
        w := t.Writer()

    for i := 0; i < 100; i++ {
        seq := w.Reserve(1)
        ring[seq&mask].Value = i
        w.Commit(seq)
    }
}

```

## Projects using Turbine

1. [CASM](https://github.com/lthibault/casm)
