# turbine

[![Godoc Reference](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/lthibault/turbine/pkg) [![Go Report Card](https://goreportcard.com/badge/github.com/lthibault/turbine?style=flat-square)](https://goreportcard.com/report/github.com/lthibault/turbine)

High-performance alternative to channels with pipelining

## Overview

Turbine is a simplified version of the now-unsupported [go-disruptor](https://github.com/smartystreets/go-disruptor).  It fixes the race conditions exhibited in go-disruptor and presents a cleaner API.

While Turbine performs better than channels in most cases, it's primary purpose is to present
a concurrent-pipeline messaging pattern.

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
    "context"
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

    go func() {
        w := t.Writer()

        for i := 0; i < 100; i++ {
            seq := w.Reserve(1)
            ring[seq&mask].Value = i
            w.Commit(seq)
        }
    }()

    <-context.Background().Done()
}

```

## Projects using Turbine

1. [CASM](https://github.com/lthibault/casm)
