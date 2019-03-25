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
		ring[seq&mask].Value *= 10
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

	var i int
	for {
		seq := w.Reserve(1)
		ring[seq&mask].Value = i
		w.Commit(seq)
		i++
	}
}
